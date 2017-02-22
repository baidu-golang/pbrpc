# pbrpc

pbrpc是一种基于TCP协议的二进制高性能RPC通信协议实现。它以Protobuf作为基本的数据交换格式。完全兼容jprotobuf-rpc-socket: https://github.com/jhunters/Jprotobuf-rpc-socket

### features:

- 内置连接池，具备更高的性能，低延迟 QPS: 5w+
- 支持自动重连功能
- 支持附件发送
- 支持超时功能
- 压缩功能，支持GZip与Snappy
- 集成内置HTTP管理功能[TODO]
- Client支持Ha的负载均衡功能[TODO]
  ​

### 依赖三方库

1. ##### golang-protobuf 针对golang开发支持google protocol  buffer库, 获取方式如下

   ##### go get github.com/golang/protobuf

2. glog 日志库, 获取方式如下

   ##### go get github.com/golang/glog

3. go-commons-pool 连接池库，获取方式如下

   ##### go get github.com/jolestar/go-commons-pool

4. 单元测试工具类 testify，获取方式如下

   ##### go get github.com/stretchr/testify/

5. link Go语言网络层脚手架，获取方式如下

   ##### go get github.com/funny/link

6. Snappy压缩类库，获取方式如下

   ##### go get github.com/golang/snappy

### Demo示例

#### 开发RPC服务端

1. 首先需要实现Service接口

   ```go
   type Service interface {
   	/*
   	   RPC service call back method.
   	   message : parameter in from RPC client or 'nil' if has no parameter
   	   attachment : attachment content from RPC client or 'nil' if has no attachment
   	   logId : with a int64 type log sequence id from client or 'nil if has no logId
   	   return:
   	   [0] message return back to RPC client or 'nil' if need not return method response
   	   [1] attachment return back to RPC client or 'nil' if need not return attachemnt
   	   [2] return with any error or 'nil' represents success
   	*/
   	DoService(message proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error)
   	GetServiceName() string
   	GetMethodName() string
       // 获得参数类型，pb会反序化到这个对象
   	NewParameter() proto.Message
   }
   ```

   ​

2. Service接口实现示例如下

```go
type SimpleService struct {
	serviceName string
	methodName  string
}

func NewSimpleService(serviceName, methodName string) *SimpleService {
	ret := SimpleService{serviceName, methodName}
	return &ret
}

func (ss *SimpleService) DoService(msg proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error) {
	var ret = "hello "

	if msg != nil {
		var name *string = nil

		m := msg.(*DataMessage)
		name = m.Name

		if len(*name) == 0 {
			ret = ret + "veryone"
		} else {
			ret = ret + *name
		}
	}
	dm := DataMessage{}
	dm.Name = proto.String(ret)
	return &dm, []byte{1, 5, 9}, nil

}

func (ss *SimpleService) GetServiceName() string {
	return ss.serviceName
}

func (ss *SimpleService) GetMethodName() string {
	return ss.methodName
}

func (ss *SimpleService) NewParameter() proto.Message {
	ret := DataMessage{}
	return &ret
}
```



3. 创建RPC服务并注册该实现接口

   ```go
   func main() {
       serverMeta := pbrpc.ServerMeta{}
   	serverMeta.Host = nil
   	serverMeta.Port = 8122
   	rpcServer := pbrpc.NewTpcServer(&serverMeta)

   	ss := NewSimpleService("echoService", "echo")

   	rpcServer.Register(ss)

   	// start server and block here
   	err := rpcServer.StartAndBlock()

   	if err != nil {
   		pbrpc.Error(err)
   		os.Exit(-1)
   	}
   }
   ```

   至此RPC已经开发完成，运行上面代码，就可以发布完成.

4. DataMessage对象定义如下

```go
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



#### 开发RPC客户端

```go
    // 创建链接(本示例使用连接池方式)
	url := pbrpc.URL{}
	url.SetHost(host).SetPort(port)
    timeout := time.Second * 5
   
    connection, err := pbrpc.NewDefaultTCPConnectionPool(url, &timeout)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
    defer connection.Close()

    // 创建client
    rpcClient, err := pbrpc.NewRpcCient(connection)
    if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
    // 调用RPC
	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := pbrpc.NewRpcInvocation(&serviceName, &methodName)
   
    // 指定压缩算法
    rpcInvocation.CompressType = proto.Int32(pbrpc.COMPRESS_GZIP)

	message := "say hello from xiemalin中文测试"
	dm := DataMessage{&message}

	rpcInvocation.SetParameterIn(&dm)
	rpcInvocation.LogId = proto.Int64(1)

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

更多使用示例参见demo