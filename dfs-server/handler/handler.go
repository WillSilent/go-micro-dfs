package handler

import (
	"context"
	"fmt"
	pb "go-micro-dfs/dfs-server/proto"
	"go-micro-dfs/dfs-server/util"
	evMsg "go-micro-dfs/service/event"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go-micro.dev/v4"
	"go-micro.dev/v4/util/log"
)
type DfsSrv struct {
	NamePub micro.Publisher
	DataPub micro.Publisher
}

type MsgIterface interface {}

// send events using the publisher
func sendEv(msg MsgIterface, p micro.Publisher) {
	t := time.NewTicker(time.Second)

	for range t.C {
		log.Logf("publishing %+v\n", msg)

		// publish an event
		if err := p.Publish(context.TODO(), msg); err != nil {
			log.Logf("error publishing: %v", err)
		}
	}
}

func (s *DfsSrv) Upload(ctx context.Context, in *pb.Args, out *pb.Result) error {
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

	// 2. 往blocker中发数据，pub-namenode，pub-datanode
	// create new event
	event := &evMsg.UpdateNameNodeEvent{
		MethodName: "Add",
		FileSha1: fileSha1,
		FileName: fileInfo.Name(),
		FileSize: fileInfo.Size(),
		ChunkNum: int32(math.Ceil(float64(fileInfo.Size()) / (64 * 1024 * 1024))),
		AddTime: time.Now().Format("2006-01-02 15:04"),
	}
	go sendEv(event, s.NamePub)

	// 3. 将文件进行分片，并放入队列中，并将文件写入到分配的datanode中，完成后，往 namenode发送一条请求，往redis中写入数据
	fileBlockPaths := splitFile(in.FilePath)
	for _, fileblockPath := range fileBlockPaths{
		for j := 0; j < 3; j++ {
			
		// 	//TODO:怎么去ip地址，保证是分块是相对散列的
			port := 2021 + j;
			SftpIPAddr := "192.168.246.100:"+strconv.Itoa(port)

			// 往blocker中发数据，pub-namenode，pub-datanode
			// create new event
			ev1 := &evMsg.UploadFile2SFTPEvent{
				FileBlockPath: fileblockPath,
				FileSha1: fileSha1,
				Replica: int32(j+1),
				SftpIPAddr: SftpIPAddr,
			}
			go sendEv(ev1, s.DataPub)

			ev2 := &evMsg.UpdateNameNodeEvent{
				MethodName: "Update",
				FileSha1: fileSha1,
				FileName: filepath.Base(fileblockPath),
				SftpIPAdr: SftpIPAddr,
				Replica: int32(j+1),
				UpdateTime: time.Now().Format("2006-01-02 15:04"),
			}
			go sendEv(ev2, s.NamePub)
		}

		log.Logf("File: %s uploaded Successfully!", fileblockPath)
	}

	// TODO: 4.删除临时文件

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
		ofile := fmt.Sprintf("./tmp/%s-%d.part", fileInfo.Name(), i)
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
