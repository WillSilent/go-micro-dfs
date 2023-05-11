package main

import (
	"fmt"
	redisConn "go-micro-dfs/namenode/db"
	"go-micro-dfs/namenode/handler"
	pb "go-micro-dfs/namenode/proto"

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
		micro.Name("dfs.NameNode"),
		micro.Version("v1.0"),
	)
	// 使用服务注册插件
	service.Init()

	// 初始化Redis连接池
	node := handler.NameNode{
		Pool: redisConn.RedisPool(),
	}

	//注册subscriber
	if err := broker.Connect(); err != nil {
		log.Error("Broker Connect Error, ", err)
	}
	//注册namenode
	_, err := broker.Subscribe("dfs.topic.namenode", eventHandler)
	if err != nil {
		log.Error("failed to Subscribe to a broker: ", err)
	}

	//rpc调用
	err = pb.RegisterNameNodeHandler(service.Server(), &node)
	
	if err != nil {
		fmt.Println("failed to register a handler: ", err)
	}

	if err = service.Run(); err != nil {
		fmt.Println("failed to run a service: ", err)
	}
}

