package handler

import (
	"context"
	"fmt"
	"go-micro-dfs/dfs-server/util"
	pbDataNode "go-micro-dfs/service/datanode"
	pbDfsSrv "go-micro-dfs/service/dfsSrv"
	pbNameNode "go-micro-dfs/service/namenode"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/util/log"
)
type DfsSrv struct {
}

// send events using the publisher
func sendEv(topic string, msg *broker.Message) {
		
	if err := broker.Publish(topic, msg); err != nil {
		log.Fatal("Broker publish error:", err)
	} else {
		log.Infof("Broker Publish a message to topic: %s.", topic)
	}
	
}

// send request using grpc
func getRpcNameNodeClient() pbNameNode.NameNodeService{
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	service := micro.NewService(micro.Registry(reg), micro.Name("Dfs.Server.client"))
	service.Init()
	return pbNameNode.NewNameNodeService("dfs.NameNode", service.Client())
}

func getRpcDataNodeClient() pbDataNode.DataNodeService{
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	service := micro.NewService(micro.Registry(reg), micro.Name("Dfs.Server.client"))
	service.Init()
	return pbDataNode.NewDataNodeService("dfs.DataNode", service.Client())
}

var opts client.CallOption = func(co *client.CallOptions) {
	co.RequestTimeout = time.Second * 30
	co.DialTimeout = time.Second * 30
}

func (s *DfsSrv) Upload(ctx context.Context, in *pbDfsSrv.Args, out *pbDfsSrv.Result) error {
	// 1. 得到文件的元数据信息，往 namenode 发送请求，将文件的元数据信息存入到redis中 （这个也可以远程通过调用namenode接口来获取文件的元数据）
	file, err := os.Open(in.FilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Logf("Get a file: %s\n", in.FilePath)
	
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	
	fileSha1 := util.FileSha1(file)

	addRsp, err := getRpcNameNodeClient().AddFileMetaData(context.TODO(),&pbNameNode.AddReq{
		FileSha1: fileSha1,
		FileName: fileInfo.Name(),
		FileSize: fileInfo.Size(),
		ChunkNum: int32(math.Ceil(float64(fileInfo.Size()) / (64 * 1024 * 1024))),
		AddTime: time.Now().Format("2006-01-02 15:04"),
	}, opts)

	if err != nil {
		log.Logf("failed to Add fileInfo to Redis: %s - %s", addRsp, err)
		return err
	}

	// 2. 往blocker中发数据，pub-namenode，pub-datanode	
	// create new event
	/** Add a file to redis
	event := &evMsg.UpdateNameNodeEvent{
		MethodName: "Add",
		FileSha1: fileSha1,
		FileName: fileInfo.Name(),
		FileSize: fileInfo.Size(),
		ChunkNum: int32(math.Ceil(float64(fileInfo.Size()) / (64 * 1024 * 1024))),
		AddTime: time.Now().Format("2006-01-02 15:04"),
	}
	msgBody, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	msg := &broker.Message{
		Header: map[string]string {
			"id": uuid.NewUUID().String(),
		},
		Body: msgBody,
	}
	sendEv("dfs.topic.namenode", msg)
	**/

	// 3. 将文件进行分片，并放入队列中，并将文件写入到分配的datanode中，完成后，往 namenode发送一条请求，往redis中写入数据
	fileBlockPaths := splitFile(in.FilePath)
	
	rand.Seed(time.Now().UnixNano())
	
	for _, fileblockPath := range fileBlockPaths{
		j := rand.Intn(3)
		SftpIPAddr := "34.145.247.183:22"
		switch j {
      case 0: SftpIPAddr = "34.145.247.183:22"
      case 1: SftpIPAddr = "34.125.193.47:22"
      case 2 : SftpIPAddr = "34.16.141.143:22"
      default: SftpIPAddr = "34.145.247.183:22"  
   }
		//TODO:怎么去ip地址，保证是分块是相对散列的
		// 34.125.193.47
		// 34.16.141.143
		// 34.145.247.183
	
		// 往blocker中发数据，pub-namenode，pub-datanode
		// create new event
		/** upload file
		ev1 := &evMsg.UploadFile2SFTPEvent{
			FileBlockPath: fileblockPath,
			FileSha1: fileSha1,
			SftpIPAddr: SftpIPAddr,
		}

		ev1MsgBody, err := json.Marshal(ev1)
		if err != nil {
			log.Fatal(err.Error())
			return nil
		}

		ev1Msg := &broker.Message{
			Header: map[string]string {
				"id": uuid.NewUUID().String(),
			},
			Body: ev1MsgBody,
		}

		sendEv("dfs.topic.datanode", ev1Msg)
		**/
		uploadRsp, err := getRpcDataNodeClient().UploadFileBlock(context.TODO(),&pbDataNode.UploadArgs{
			FileBlockPath: fileblockPath,
			FileSha1: fileSha1,
			SftpIPAddr: SftpIPAddr,
		}, opts)
	
		if err != nil {
			log.Logf("failed to Upload File to SFTP: %s - %s", uploadRsp, err)
			return err
		}
		// 往redis回写数据
		/**
		ev2 := &evMsg.UpdateNameNodeEvent{
			MethodName: "Update",
			FileSha1: fileSha1,
			FileName: filepath.Base(fileblockPath),
			SftpIPAdr: SftpIPAddr,
			Replica: 1,
			UpdateTime: time.Now().Format("2006-01-02 15:04"),
		}
		ev2MsgBody, err := json.Marshal(ev2)
		if err != nil {
			log.Fatal(err.Error())
			return nil
		}

		ev2Msg := &broker.Message{
			Header: map[string]string {
				"id": uuid.NewUUID().String(),
			},
			Body: ev2MsgBody,
		}
		sendEv("dfs.topic.namenode", ev2Msg)
		**/
		updateRsp, err := getRpcNameNodeClient().UpdateFileBlockMetaData(context.TODO(),&pbNameNode.UpdateReq{
			FileSha1: fileSha1,
			FileName: filepath.Base(fileblockPath),
			SftpIPAdr: SftpIPAddr,
			Replica: 1,
			UpdateTime: time.Now().Format("2006-01-02 15:04"),
		}, opts)
	
		if err != nil {
			log.Logf("failed to Add fileInfo to Redis: %s - %s", updateRsp, err)
			return err
		}

		fmt.Printf("File: %s Published!\n", fileblockPath)
	}

	// TODO: 4.删除临时文件
	out.Code = 200
	out.Message = "Upload Finshed!";
	return nil
}

// 文件切分函数
func splitFile(infile string) []string{
	if infile == "" {
		log.Fatal("请输入正确的文件名")
		return nil
	}
	
	fileInfo, err := os.Stat(infile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("File doesn't exist!")
			return nil
		}
		log.Fatal(err)
		return nil
	}

	var chunkSize int64 = 64 * 1024 * 1024

	num := int(math.Ceil(float64(fileInfo.Size()) / (64 * 1024 * 1024)))
	filepaths := make([]string, num)

	fi, err := os.OpenFile(infile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return filepaths
	}
	fmt.Printf("The file will be splited into %d pieces.\n", num)
	
	b := make([]byte, chunkSize)
	var i int64 = 1
	for ; i <= int64(num); i++ {
		fi.Seek((i-1) * chunkSize, 0)
		if len(b) > int(fileInfo.Size()-(i-1) * chunkSize) {
			b = make([]byte, fileInfo.Size()-(i-1) * chunkSize)
		}
		fi.Read(b)

		//TODO: filename need to be modified.
		ofile := fmt.Sprintf("/Users/will/Desktop/git-code/go-micro-dfs/test/tmp/%s-%d.part", fileInfo.Name(), i)
		filepaths[i-1] = ofile 
		
		fmt.Printf("Create: %s\n", ofile)
		f, err := os.OpenFile(ofile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		f.Write(b)
		f.Close()
	}
	fi.Close()
	fmt.Println("Split Finished!")

	return filepaths
}
