package main

import (
	"encoding/json"
	"fmt"
	redisConn "go-micro-dfs/namenode/db"
	"go-micro-dfs/namenode/handler"
	pb "go-micro-dfs/namenode/proto"
	evMsg "go-micro-dfs/service/event"
	"strconv"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/gomodule/redigo/redis"
	"go-micro.dev/v4"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/util/log"
)

var Pool *redis.Pool

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
	Pool = redisConn.RedisPool()

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
	node := handler.NameNode{
		Pool: Pool,
	}
	err = pb.RegisterNameNodeHandler(service.Server(), &node)
	
	if err != nil {
		fmt.Println("failed to register a handler: ", err)
	}

	if err = service.Run(); err != nil {
		fmt.Println("failed to run a service: ", err)
	}
}

func eventHandler(b broker.Event) error{
	var req *evMsg.UpdateNameNodeEvent
	if err := json.Unmarshal(b.Message().Body, &req); err != nil {
			return err
	}
	fmt.Printf("Msg: %s:%s\n", req.FileName, req.MethodName)
	//doUpdate or doAdd
	if req.MethodName == "Add" {
		doAdd(req)
	}

	if req.MethodName == "Update" {
		doUpdate(req)
	}
	return nil
}

func doAdd(req *evMsg.UpdateNameNodeEvent) error{
	// 调用redis的连接池，并返回上传结果
	rConn := Pool.Get()

	if rConn == nil {
		log.Fatal("Error: Cannot not get a connection from redis.")
		return nil
	}

	defer rConn.Close()
	//将初始信息写入redis缓存
	rConn.Do("HSET", req.FileSha1, "filename", req.FileName)
	rConn.Do("HSET", req.FileSha1, "filesize", req.FileSize)
	rConn.Do("HSET", req.FileSha1, "chunkcount", req.ChunkNum)
	rConn.Do("HSET", req.FileSha1, "addAt", req.AddTime)
	log.Infof("Successfully Add a File Meta!")
	return nil
}

func doUpdate(req *evMsg.UpdateNameNodeEvent) error{
	// 调用redis的连接池，并返回上传结果
	rConn := Pool.Get()

	if rConn == nil {
		log.Fatal("Error: Cannot not get a connection from redis.")
		return nil
	}

	defer rConn.Close()
	//更新: 还要加一个update_time
	rConn.Do("HSET", req.FileSha1, req.FileName + "_replica_"+strconv.Itoa(int(req.Replica)), req.SftpIPAdr)
	rConn.Do("HSET", req.FileSha1, "updateAt", req.UpdateTime)
	log.Infof("Successfully Update file block data in the redis !")
	return nil
}

