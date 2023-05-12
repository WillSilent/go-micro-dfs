package main

import (
	"context"
	"fmt"
	pb "go-micro-dfs/dfs-server/proto"
	"time"

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

	service := micro.NewService(micro.Registry(reg), micro.Name("client"))
	service.Init()

	fmt.Println("Begin to upload file.....")
	t1 := time.Now()
	rsp, err := pb.NewDfsSrvService("dfs.Server", service.Client()).Upload(context.TODO(), &pb.Args{FilePath: "C:/Users/will9/OneDrive/Desktop/go-micro-dfs/test/test.mp4"})
	if err != nil {
		fmt.Println("failed to new DfsSrvService: ", err)
	}
	t2 := time.Now()
	fmt.Println(rsp.Code, rsp.Message)
	fmt.Printf("Time Spend:%s.", t2.Sub(t1))
}