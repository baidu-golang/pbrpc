<!--
 * @Author: Malin Xie
 * @Description: 
 * @Date: 2021-08-23 16:01:03
-->
<h1 align="center">baidurpc</h1>

<p align="center">
baidurpc是一种基于TCP协议的二进制高性能RPC通信协议实现。它以Protobuf作为基本的数据交换格式。
本版本基于golang实现.完全兼容jprotobuf-rpc-socket: https://github.com/Baidu-ecom/Jprotobuf-rpc-socket
</p>



### 同步发布http rpc服务
要同步发布http rpc 服务，部署非常简单，直接调用 EnableHttp 方法即可。会与rpc复用同一端口。

```go
rpcServer.EnableHttp()
```

```go
	serverMeta := baidurpc.ServerMeta{}
	serverMeta.Host = nil
	serverMeta.Port = Int(*port)
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	echoService := new(EchoService)

	rpcServer.Register(echoService)
	// 开启http rpc服务
	rpcServer.EnableHttp()
	// 启动RPC服务
	err := rpcServer.Start()

	if err != nil {
		baidurpc.Error(err)
	}
```

### 相关事项说明
1. 发布http rpc 服务， 需要struct 增加json tag字段说明
```go
type DataMessage struct {
	Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
}
```  
2. http 协议统一使用POST方法进行发布。 url path 协同 http://host:port/rpc/service_name/method_name
3. http请求内容
   header 定义 "Content-Type"  "application/json;charset=utf-8"
   POST body直接为json结构
4. http返回内容, json结果.  
   <pre>
   errno错误码： 0表示成功，其它则为错误
   message: 错误信息，errno非0时，会设置
   data：返回内容
   </pre>
```json   
   {
	   "errno" : 0,
	   "message" : "",
	   "data" : {}
   }
```