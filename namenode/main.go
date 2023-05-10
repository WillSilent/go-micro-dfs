package main

import (
	"fmt"
	redisConn "go-micro-dfs/namenode/db"
	"go-micro-dfs/namenode/handler"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

// 1. 连接好redis
// 2. 把文件元数据信息写入redis服务器中，同时响应datanode的写回请求
func main() {
	// consul 服务地址按照实际情况填写
	reg := consul.NewRegistry(registry.Addrs("192.168.246.100:8500"))

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
	err := micro.RegisterSubscriber("dfs.topic.namenode", service.Server(), node.Handler)
	
	//rpc调用
	//err := pb.RegisterNameNodeHandler(service.Server(), &node)
	
	if err != nil {
		fmt.Println("failed to register a handler: ", err)
	}

	if err = service.Run(); err != nil {
		fmt.Println("failed to run a service: ", err)
	}
}

