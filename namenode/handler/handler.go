package handler

import (
	"context"
	pb "go-micro-dfs/namenode/proto"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"go-micro.dev/v4/logger"
)

type NameNode struct{
	Pool *redis.Pool
}

func (n *NameNode) AddFileMetaData(ctx context.Context, req *pb.AddReq, rsp *pb.Result) error {
	// 调用redis的连接池，并返回上传结果
	rConn := n.Pool.Get()

	if rConn == nil {
		logger.Fatal("Error: Cannot not get a connection from redis.")
		rsp.Code = 500
		rsp.Message = "Cannot not get a connection from redis."
		return nil
	}

	defer rConn.Close()

	//将初始信息写入redis缓存
	rConn.Do("HSET", req.FileSha1, "filename", req.FileName)
	rConn.Do("HSET", req.FileSha1, "filesize", req.FileSize)
	rConn.Do("HSET", req.FileSha1, "chunkcount", req.ChunkNum)
	rConn.Do("HSET", req.FileSha1, "addAt", req.AddTime)
	logger.Info("Successfully Add a File Meta!")
	return nil
}

func (n *NameNode) UpdateFileBlockMetaData(ctx context.Context, req *pb.UpdateReq, rsp *pb.Result) error {
	// 调用redis的连接池，并返回上传结果
	rConn := n.Pool.Get()

	if rConn == nil {
		logger.Fatal("Error: Cannot not get a connection from redis.")
		rsp.Code = 500
		rsp.Message = "Cannot not get a connection from redis."
		return nil
	}

	defer rConn.Close()

	//更新: 还要加一个update_time
	rConn.Do("HSET", req.FileSha1, req.FileName + "_replica_"+strconv.Itoa(int(req.Replica)), req.SftpIPAdr)
	rConn.Do("HSET", req.FileSha1, "updateAt", req.UpdateTime)
	logger.Info("Successfully Update file block data in the redis !")
	return nil

}