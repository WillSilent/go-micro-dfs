package main

import (
	"fmt"
	"go-micro-dfs/datanode/handler"
	pb "go-micro-dfs/datanode/proto"
	sftpUtil "go-micro-dfs/datanode/util"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/util/log"
)


func eventHandler(b broker.Event) error{
	msg := string(b.Message().Body)
	logger.Infof("message: %s, header: %s", msg, b.Message().Header)
	return nil
}

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
	node := handler.DataNode {
		ConnManager: sftpUtil.NewSFTPConnectionManager(10),
	}

	//注册subscriber
	// err := micro.RegisterSubscriber("dfs.topic.datanode", service.Server(), new(Sub))
	// if err != nil {
	// 	fmt.Println("failed to register a subcriber: ", err)
	// 	return
	// }
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
	err = pb.RegisterDataNodeHandler(service.Server(), &node)
	if err != nil {
		fmt.Println("failed to register a handler: ", err)
		return
	}

	if err = service.Run(); err != nil {
		fmt.Println("failed to run a service: ", err)
	}
}