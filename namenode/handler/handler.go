package main

import (
	"context"
	pb "go-micro-dfs/namenode/proto"
)

type NameNode struct{}

func (n *NameNode) AddFileMetaData(ctx context.Context, req *pb.AddReq, rsp *pb.Result) error {
	// 调用redis的连接池，并返回上传结果
	return nil
}

func (n *NameNode) UpdateFileBlockMetaData(ctx context.Context, req *pb.UpdateReq, rsp *pb.Result) error {
	// 调用redis的连接池，并返回上传结果
	return nil
}