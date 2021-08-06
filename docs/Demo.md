<h1 align="center">baidurpc</h1>

<p align="center">
baidurpc是一种基于TCP协议的二进制高性能RPC通信协议实现。它以Protobuf作为基本的数据交换格式。
本版本基于golang实现.完全兼容jprotobuf-rpc-socket: https://github.com/Baidu-ecom/Jprotobuf-rpc-socket
</p>

### 更多特性使用介绍

#### 开发RPC服务端

1. 定义PB对象
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
2. 定义一个对象以及方法，用于发布服务

   ```go

	type EchoService struct {
	}

	// Echo  test publish method with return type has context argument
	// 方法要求
	// 参数个数必须为1个或2个， 第一个类型必须为 context.Context 
	// 第二个类型必须是实现 proto.Message接口(如果是无参，可以省略)
	// 返回个数可以为1个或2个  第一个类型必须是实现 proto.Message接口 
	// 第2个参数为可选。 当使用时，必须为 context.Context类型
	func (rpc *EchoService) Echo(c context.Context, in *DataMessage) (*DataMessage, context.Context) {
		var ret = "hello "

		// if receive with attachement
		attachement := baidurpc.Attachement(c)
		fmt.Println(c)

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

以下都是合法的定义方法
1. Echo(c context.Context, in *DataMessage) (*DataMessage, context.Context)
2. Echo(c context.Context) (*DataMessage, context.Context)
3. Echo(c context.Context) (*DataMessage)


2. 指定发布端口，把EchoService发布成RPC服务

```go
	serverMeta := baidurpc.ServerMeta{}
	serverMeta.Host = nil
	serverMeta.Port = Int(*port)
	// set chunk size this will open server chunk package by specified size 
	// serverMeta.ChunkSize = 1024 // 1k
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	echoService := new(EchoService)

    // mapping可选，如果需要映射成新的function名称时使用
	mapping := make(map[string]string)
	mapping["Echo"] = "echo"
	// 第一个参数 "echoService" 为空时，则会使用 EchoService的struct 的type name
	rpcServer.RegisterNameWithMethodMapping("echoService", echoService, mapping)
	// 最简注册方式 rpcServer.Register(echoService)

	// 启动RPC服务
	err := rpcServer.StartAndBlock()

	if err != nil {
		baidurpc.Error(err)
		os.Exit(-1)
	}
```

   至此RPC已经开发完成，运行上面代码，就可以发布完成.

#### 开发启验证功能

实现 AuthService 接口

```go
type StringMatchAuthService struct {
}

// Authenticate
func (as *StringMatchAuthService) Authenticate(service, name string, authToken []byte) bool {
	if authToken == nil {
		return false
	}
	return strings.Compare(AUTH_TOKEN, string(authToken)) == 0
}

```
设置到service对象
```go

// ...
rpcServer := baidurpc.NewTpcServer(&serverMeta)
rpcServer.SetAuthService(new(StringMatchAuthService))

```

#### 设置trace功能

实现 TraceService 接口

```go
type AddOneTraceService struct {
}

// Trace
func (as *AddOneTraceService) Trace(service, name string, traceInfo *baidurpc.TraceInfo) *baidurpc.TraceInfo {
	*traceInfo.SpanId++
	*traceInfo.TraceId++
	*traceInfo.ParentSpanId++
	return traceInfo
}

```

设置到service对象
```go

// ...
rpcServer := baidurpc.NewTpcServer(&serverMeta)
rpcServer.SetTraceService(new(AddOneTraceService))

```

### 开发RPC客户端

```go
    // 创建链接(本示例使用连接池方式)
	url := baidurpc.URL{}
	url.SetHost(host).SetPort(port)
    timeout := time.Second * 5
   
    // 创建连接 
    connection, err := baidurpc.NewDefaultTCPConnectionPool(url, &timeout)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
    defer connection.Close()

    // 创建client
    rpcClient, err := baidurpc.NewRpcCient(connection)
    if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer rpcClient.Close()
    // 调用RPC
	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	message := "say hello from xiemalin中文测试"
	dm := DataMessage{&message}

	rpcInvocation.SetParameterIn(&dm)
	rpcInvocation.LogId = proto.Int64(1)

	// 可选， 设置logid 与  附件 
	// rpcInvocation.LogId = proto.Int64(1)
	// rpcInvocation.Attachment = []byte("this is attachement contenet")

	parameterOut := DataMessage{}

	response, err := rpcClient.SendRpcRequest(rpcInvocation, &parameterOut)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if response == nil {
		fmt.Println("Reponse is nil")
		return
	}
```

### 设置调用超时

```go
	// baidurpc的超时控制功能使用了 时间轮 timewheel功能 https://github.com/jhunters/timewheel
	// 可以在初始化Client时设置
	timewheelInterval := 1 * time.Second
	var timewheelSlot uint16 = 300 
	rpcClient, err := baidurpc.NewRpcCientWithTimeWheelSetting(connection, timewheelInterval, timewheelSlot)

    // 调用时，设置超时功能
	response, err := rpcClient.SendRpcRequestWithTimeout(100*time.Millisecond, rpcInvocation, &parameterOut)
	// 如果发生超时， 返回的错误码为 62

```

### 设置验证
```go
    // 调用RPC
	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)
	// set auth token
	rpcInvocation.AuthenticateData = []byte("AUTH_TOKEN")
    // 调用时，设置超时功能
	response, err := rpcClient.SendRpcRequestWithTimeout(100*time.Millisecond, rpcInvocation, &parameterOut)
	// 如果发生超时， 返回的错误码为 62

```

### 设置分包chunk功能
```go
    // 调用RPC
	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)
	// 设置分包大小(byte)
	rpcInvocation.ChunkSize = 1024 //1k
    // 调用时，设置超时功能
	response, err := rpcClient.SendRpcRequestWithTimeout(100*time.Millisecond, rpcInvocation, &parameterOut)
	// 如果发生超时， 返回的错误码为 62

```

### 设置Trace功能
```go
    // 调用RPC
	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)
	// 设置trace信息
	rpcInvocation.TraceId = 10
	rpcInvocation.SpanId = 11
	rpcInvocation.ParentSpanId = 12
	rpcInvocation.RpcRequestMetaExt = map[string]string{"key1": "value1"}
    // 调用时，设置超时功能
	response, err := rpcClient.SendRpcRequestWithTimeout(100*time.Millisecond, rpcInvocation, &parameterOut)
	// 如果发生超时， 返回的错误码为 62

	// 获取服务端返回的trace信息
	response.GetTraceId()
	response.GetParentSpanId()
	response.GetParentSpanId()
	response.GetRpcRequestMetaExt()

```


### 开发Ha RPC客户端

```go
 	urls := []baidurpc.URL{{Host: host, Port: &errPort}, {Host: host, Port: port}}

	connections, err := baidurpc.NewBatchTCPConnection(urls, timeout)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer baidurpc.CloseBatchConnection(connections)

	haClient, err := baidurpc.NewHaRpcCient(connections)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	message := "say hello from xiemalin中文测试"
	dm := DataMessage{&message}

	rpcInvocation.SetParameterIn(&dm)
	rpcInvocation.LogId = proto.Int64(1)
	rpcInvocation.Attachment = []byte("hello world")

	parameterOut := DataMessage{}

	response, err := haClient.SendRpcRequest(rpcInvocation, &parameterOut)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if response == nil {
		fmt.Println("Reponse is nil")
		return
	}

	fmt.Println("attachement", response.Attachment)

```
