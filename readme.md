## 使用 go-micro@v4 来搭建 dfs

1. 需要先安裝 consul，以 consul 作为注册中心

2. 所需组件：

3. 后续可以增加一个 broker/使用 kafka 或者 rabbitmq

记录几个值得看的博客：
go-micro 框架：

- https://zhuanlan.zhihu.com/p/372796932
- https://www.cnblogs.com/xiangxiaolin/p/12820837.html
- https://medium.com/@dche423/micro-in-action-getting-start-cn-99c870e078f

官方例子：https://github.com/go-micro/examples
官方 generator： https://github.com/go-micro/generator

可以参考一下的项目：https://github.com/xbox1994/go-micro-example
https://juejin.cn/post/7006228288708444173

```bash
  #v4版本及以上
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest

  #v3版本
  go get github.com/micro/micro/v3/cmd/protoc-gen-micro
```

### 搭建 sftp 服务器

https://zhuanlan.zhihu.com/p/494150542
