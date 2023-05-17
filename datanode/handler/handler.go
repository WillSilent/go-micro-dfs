package handler

import (
	"context"
	sftpUtil "go-micro-dfs/datanode/util"
	pbDataNode "go-micro-dfs/service/datanode"

	"go-micro.dev/v4/logger"
)

type DataNode struct {
	ConnManager *sftpUtil.SFTPConnectionManager
}

func (d *DataNode) UploadFileBlock(ctx context.Context, req *pbDataNode.UploadArgs, rsp *pbDataNode.DnodeResult) error {
	pool, err := d.ConnManager.GetPool(req.SftpIPAddr, "sftpuser", "123456")
	if err != nil {
		logger.Fatal("Error:", err)
		rsp.Code = 500
		rsp.Message = "Cannot not get a connection from sftp server."
		return nil
	}

	client := pool.Get()
	defer pool.Put(client)

	// 在此处执行SFTP操作，例如client.client.ReadDir("/path/to/directory")
	sftpUtil.UploadFile(client, req.FileBlockPath, "/data")
	logger.Info("Error:", err)
	return nil
}
