package main

import (
	"context"
	"fmt"
	pbDfsSrv "go-micro-dfs/service/dfsSrv"
	"time"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
)

var opts client.CallOption = func(co *client.CallOptions) {
	co.RequestTimeout = time.Second * 30
	co.DialTimeout = time.Second * 30
}

func main() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	service := micro.NewService(micro.Registry(reg), micro.Name("greeter.client"))
	service.Init()

	fmt.Println("Begin to upload file.....")
	t1 := time.Now()
	rsp, err := pbDfsSrv.NewDfsSrvService("dfs.Server", 
									service.Client(),
							).Upload(context.TODO(), 
												&pbDfsSrv.Args{FilePath: "/Users/will/Desktop/git-code/go-micro-dfs/test/test_500M.mp4"}, 
												opts,
							)
	if err != nil {
		fmt.Println("failed to new greeter service: ", err)
	}
	t2 := time.Now()
	fmt.Println(rsp.Code, rsp.Message)
	fmt.Printf("Time Spend:%s.", t2.Sub(t1))
}