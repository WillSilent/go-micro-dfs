package main

import (
	"fmt"
	"go-micro-dfs/datanode/handler"
	sftpUtil "go-micro-dfs/datanode/util"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

// 1.建立sftp连接池
// 2.把文件上传至sftp服务器
func main() {
	// consul 服务地址按照实际情况填写
	reg := consul.NewRegistry(registry.Addrs("192.168.246.100:8500"))

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
	err := micro.RegisterSubscriber("dfs.topic.datanode", service.Server(), node.Handler)
	
	//err := pb.RegisterDataNodeHandler(service.Server(), &node)
	if err != nil {
		fmt.Println("failed to register a handler: ", err)
	}

	if err = service.Run(); err != nil {
		fmt.Println("failed to run a service: ", err)
	}
}