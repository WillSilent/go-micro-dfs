package main

import (
	"context"
	"fmt"
	pb "go-micro-dfs/hello/proto"
	"log"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	log.Printf("I received this request.")
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {

	 // consul 服务地址按照实际情况填写
	 reg := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("hello"),
		micro.Version("v1.0"),
	)
	// 使用服务注册插件
	service.Init()

	err := pb.RegisterHelloHandler(service.Server(), new(Greeter))
	if err != nil {
		fmt.Println("failed to register a handler: ", err)
	}

	if err = service.Run(); err != nil {
		fmt.Println("failed to run a service: ", err)
	}
}