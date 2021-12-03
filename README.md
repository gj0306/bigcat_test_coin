# news_coin


### 介绍：

bigcat_test_coin 是为了做一些区块儿链的共识机制的验证，并不具备成熟的商业价值。

##### 共识：
bigcat_test_coin是在拜占庭将军决策的基础上引入了**时间证明**概念，即每分钟出一个Block，并且时间也是Block的一个重要的验证参数。

出块儿节点采用的方案是，所在的在所有的线节点每一个时间周期都可以发布Block，每个Block都有着自身的权重值，节点会自动跳到权重最高的链中，这样可以做到完全去中心化的效果。

账本并没有采取utxo模型，而是account模型。

###  安装:

```go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
go install code.google.com/p/goprotobuf/proto@latest

# 安装 protoc 并加入环境变量
https://github.com/protocolbuffers/protobuf/tags
```

###  运行：

```go
执行 编译win.cmd 或 编译linux.sh 文件

```

