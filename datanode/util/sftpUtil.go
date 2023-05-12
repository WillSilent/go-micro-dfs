package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 连接池
type SFTPConnectionPool struct {
	conns chan *sftp.Client
}

func NewSFTPConnectionPool(capacity int, addr, user, password string) (*SFTPConnectionPool, error) {
	conns := make(chan *sftp.Client, capacity)

	for i := 0; i < capacity; i++ {
		client, err := connectSFTP(addr, user, password)
		if err != nil {
			return nil, err
		}
		conns <- client
	}

	return &SFTPConnectionPool{conns: conns}, nil
}

func connectSFTP(addr, user, password string) (*sftp.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return client, nil
}

func (cp *SFTPConnectionPool) Get() *sftp.Client {
	return <-cp.conns
}

func (cp *SFTPConnectionPool) Put(client *sftp.Client) {
	cp.conns <- client
}

func (cp *SFTPConnectionPool) Close() {
	for i := 0; i < cap(cp.conns); i++ {
		(<-cp.conns).Close()
	}
}

type SFTPConnectionManager struct {
	mu       sync.Mutex
	pools    sync.Map
	capacity int
}

func NewSFTPConnectionManager(capacity int) *SFTPConnectionManager {
	return &SFTPConnectionManager{capacity: capacity}
}

func (cm *SFTPConnectionManager) GetPool(serverAddr, user, password string) (*SFTPConnectionPool, error) {
	key := serverAddr + ":" + user + ":" + password
	pool, ok := cm.pools.Load(key)
	if ok {
		return pool.(*SFTPConnectionPool), nil
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	pool, ok = cm.pools.Load(key)
	if ok {
		return pool.(*SFTPConnectionPool), nil
	}

	newPool, err := NewSFTPConnectionPool(cm.capacity, serverAddr, user, password)
	if err != nil {
		return nil, err
	}

	cm.pools.Store(key, newPool)
	return newPool, nil
}

// 只获取一个连接
func Connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
	 auth   []ssh.AuthMethod
	 addr   string
	 clientConfig *ssh.ClientConfig
	 sshClient *ssh.Client
	 sftpClient *sftp.Client
	 err   error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
	 User:   user,
	 Auth:   auth,
	 Timeout:   30 * time.Second,
	 HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
	 return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
	 return nil, err
	}
	return sftpClient, nil
 }

 /**
  *  SFTP服务器操作文件
	*/
 func UploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
	 fmt.Println("os.Open error : ", localFilePath)
	 log.Fatal(err)
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
	 fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
	 log.Fatal(err)
	}
	defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
	 fmt.Println("ReadAll error : ", localFilePath)
	 log.Fatal(err)
	}
	dstFile.Write(ff)
	log.Printf(localFilePath + " copy file to remote server finished!\n")
 }

 func UploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
	 log.Fatal("read dir list fail ", err)
	}
	for _, backupDir := range localFiles {
	 localFilePath := path.Join(localPath, backupDir.Name())
	 remoteFilePath := path.Join(remotePath, backupDir.Name())
	 if backupDir.IsDir() {
		sftpClient.Mkdir(remoteFilePath)
		UploadDirectory(sftpClient, localFilePath, remoteFilePath)
	 } else {
		UploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
	 }
	}
	fmt.Println(localPath + " copy directory to remote server finished!")
 }