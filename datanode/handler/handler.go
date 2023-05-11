package handler

import (
	"context"
	pb "go-micro-dfs/datanode/proto"
	sftpUtil "go-micro-dfs/datanode/util"

	"go-micro.dev/v4/logger"
)

type DataNode struct {
	ConnManager *sftpUtil.SFTPConnectionManager
}

func (d *DataNode) UploadFileBlock(ctx context.Context, req *pb.UploadArgs, rsp *pb.Result) error {
	pool, err := d.ConnManager.GetPool(req.SftpIPAddr, "admin", "admin")
	if err != nil {
		logger.Fatal("Error:", err)
		rsp.Code = 500
		rsp.Message = "Cannot not get a connection from sftp server."
		return nil
	}

	client := pool.Get()
	defer pool.Put(client)

	// 在此处执行SFTP操作，例如client.client.ReadDir("/path/to/directory")
	sftpUtil.UploadFile(client, req.FileBlockPath, "/upload")
	logger.Info("Error:", err)
	return nil
}
