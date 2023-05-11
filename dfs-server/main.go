package main

import (
	"fmt"
	"go-micro-dfs/dfs-server/handler"
	pb "go-micro-dfs/dfs-server/proto"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

// 生成文件元数据信息，并对文件进行切分
func main() {
	// consul 服务地址按照实际情况填写
	reg := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("dfs.Server"),
		micro.Version("v1.0"),
	)
	// 使用服务注册插件
	service.Init()

	// create publisher
	srv := handler.DfsSrv{
		NamePub: micro.NewEvent("dfs.topic.namenode", service.Client()),
		DataPub: micro.NewEvent("dfs.topic.datanode", service.Client()),
	}

	err := pb.RegisterDfsSrvHandler(service.Server(), &srv)
	if err != nil {
		fmt.Println("failed to register a handler: ", err)
	}

	if err = service.Run(); err != nil {
		fmt.Println("failed to run a service: ", err)
	}
}