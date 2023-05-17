package main

import (
	"context"
	"fmt"
	pbDfsSrv "go-micro-dfs/service/dfsSrv"

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

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("client"), //name the client service
	)
	// Initialise service
	service.Init()

	//create hello service client
	client := pbDfsSrv.NewDfsSrvService("dfs.Server", service.Client())

	
	rsp, err := client.Upload(context.TODO(), &pbDfsSrv.Args{FilePath: "D:/code/github_code/dfs/test/test.mp4"})
	if err != nil {
		fmt.Println("failed to new dfsServer service: ", err)
		return
	}

	fmt.Println(rsp.Code, rsp.Message)
}