<!--
 * @Author: Malin Xie
 * @Description: 
 * @Date: 2021-08-06 13:14:23
-->
<h1 align="center">baidurpc</h1>

<p align="center">
baidurpc是一种基于TCP协议的二进制高性能RPC通信协议实现。它以Protobuf作为基本的数据交换格式。
本版本基于golang实现.完全兼容jprotobuf-rpc-socket: https://github.com/Baidu-ecom/Jprotobuf-rpc-socket
</p>


### Installing 

To start using pbrpc, install Go and run `go get`:

```sh
$ go get github.com/baidu-golang/pbrpc
```

### 定义protobuf 对象
   ```property
   message DataMessae {
	   string name = 1;
   }
	```
	用protoc工具 生成 pb go 定义文件
	protoc --go_out=. datamessage.proto


### 开发RPC服务端
要发布一个RPC服务，需要先定义一个对象以及方法，用于发布服务。 与传统的golang基础库的rpc发布非常一致。

1. 以下定义了 <b>EchoService对象</b> 并 增加 <b>Echo</b> 方法
   ```go
	type EchoService struct {
	}

	func (rpc *EchoService) Echo(c context.Context, in *DataMessage) (*DataMessage, context.Context) {
		var ret = "hello "
		if len(*in.Name) == 0 {
			ret = ret + "veryone"
		} else {
			ret = ret + *in.Name
		}
		dm := DataMessage{}
		dm.Name = proto.String(ret)
		return &dm, baidurpc.BindAttachement(context.Background(), []byte("hello")) // return with attachement
	}
   ```

2. 指定发布端口，把EchoService发布成RPC服务

```go
	serverMeta := baidurpc.ServerMeta{}
	serverMeta.Host = nil
	serverMeta.Port = Int(*port)
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	echoService := new(EchoService)

	rpcServer.Register(echoService)

	// 启动RPC服务
	err := rpcServer.Start()

	if err != nil {
		baidurpc.Error(err)
	}
```

 [查看代码](https://github.com/baidu-golang/pbrpc/blob/master/example/server_example_test.go) <br>