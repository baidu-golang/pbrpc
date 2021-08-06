<!--
 * @Author: Malin Xie
 * @Description: 
 * @Date: 2021-08-06 13:37:43
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
以下按pb2进行定义
   ```go
   //手工定义pb生成的代码, tag 格式 = protobuf:"type,order,req|opt|rep|packed,name=fieldname"
	type DataMessage struct {
		Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	}

	func (m *DataMessage) Reset()         { *m = DataMessage{} }
	func (m *DataMessage) String() string { return proto.CompactTextString(m) }
	func (*DataMessage) ProtoMessage()    {}

	func (m *DataMessage) GetName() string {
		if m.Name != nil {
			return *m.Name
		}
		return ""
	}
	```


### 开发RPC客户端
客户端开发，使用RpcClient进行调用。

```go
	host := "localhost"
	port := 1031
    // 创建链接
	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 500
	// create client by simple connection
	connection, err := baidurpc.NewTCPConnection(url, &timeout)
   
	if err != nil {
		fmt.Println(err)
		return
	}
    defer connection.Close()

    // 创建client
    rpcClient, err := baidurpc.NewRpcCient(connection)
    if err != nil {
		fmt.Println(err)
		return
	}
	defer rpcClient.Close()
    // 调用RPC
	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	message := "say hello from xiemalin中文测试"
	dm := DataMessage{&message}

	rpcInvocation.SetParameterIn(&dm)
	parameterOut := DataMessage{}

	_, err := rpcClient.SendRpcRequest(rpcInvocation, &parameterOut)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(*parameterOut.Name)
```

