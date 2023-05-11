package main

import (
	"context"
	"fmt"
	pb "go-micro-dfs/dfs-server/proto"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

func main() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	service := micro.NewService(micro.Registry(reg), micro.Name("greeter.client"))
	service.Init()
	rsp, err := pb.NewDfsSrvService("dfs.Server", service.Client()).Upload(context.TODO(), &pb.Args{FilePath: "/Users/will/Desktop/git-code/go-micro-dfs/test//test.mov"})
	if err != nil {
		fmt.Println("failed to new greeter service: ", err)
	}

	fmt.Println(rsp.Code, rsp.Message)
}