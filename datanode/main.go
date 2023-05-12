package main

import (
	"encoding/json"
	"fmt"
	"go-micro-dfs/datanode/handler"
	pb "go-micro-dfs/datanode/proto"
	sftpUtil "go-micro-dfs/datanode/util"
	evMsg "go-micro-dfs/service/event"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/util/log"
)

var cm *sftpUtil.SFTPConnectionManager

func main() {
	// consul 服务地址按照实际情况填写
	reg := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("dfs.DataNode"),
		micro.Version("v1.0"),
	)
	// 使用服务注册插件
	service.Init()

	// 初始化SFTP连接池
	cm = sftpUtil.NewSFTPConnectionManager(10)
	
	//注册subscriber
	if err := broker.Connect(); err != nil {
		log.Error("Broker Connect Error, ", err)
	}
	//注册namenode
	_, err := broker.Subscribe("dfs.topic.datanode", eventHandler)
	if err != nil {
		log.Error("failed to Subscribe to a broker: ", err)
	}

	//注册grpc handler
	node := handler.DataNode {
		ConnManager: cm,
	}
	err = pb.RegisterDataNodeHandler(service.Server(), &node)
	if err != nil {
		fmt.Println("failed to register a handler: ", err)
		return
	}

	if err = service.Run(); err != nil {
		fmt.Println("failed to run a service: ", err)
	}
}

func eventHandler(b broker.Event) error{
	// msg := string(b.Message().Body)
	// logger.Infof("message: %s, header: %s", msg, b.Message().Header)
	var req *evMsg.UploadFile2SFTPEvent
	if err := json.Unmarshal(b.Message().Body, &req); err != nil {
			return err
	}
	fmt.Printf("The file: %s will be uploaded to: %s.\n", req.FileBlockPath, req.SftpIPAddr)
	pool, err := cm.GetPool(req.SftpIPAddr, "sftpuser", "123456")
	if err != nil {
		log.Fatal("Error:", err)
		return nil
	}

	client := pool.Get()
	if client == nil {
		log.Fatal("Cannot get a client form pool:", err)
		return nil
	}
	defer pool.Put(client)

	// 在此处执行SFTP操作，例如client.client.ReadDir("/path/to/directory")
	sftpUtil.UploadFile(client, req.FileBlockPath, "/data")
	return nil
}