package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"
	pool "github.com/jolestar/go-commons-pool/v2"

	"github.com/golang/protobuf/proto"
)

var host = flag.String("host", "localhost", "If non-empty, connect to target host")
var port = flag.Int("port", 8122, "If non-empty, use this port")

var countsPerThread = flag.Int("countsPerThread", 1, "If non-empty, counts per thread to execute default is 100")
var parell = flag.Int("parell", 1, "If non-empty, concurrent threads to run")

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

func main() {

	url := baidurpc.URL{}
	url.SetHost(host).SetPort(port)
	fmt.Println(url)
	timeout := time.Second * 5
	// create client by connection pool
	config := pool.NewDefaultPoolConfig()
	config.MaxTotal = *parell
	config.MaxIdle = *parell
	config.TestOnBorrow = true

	// connection, err := baidurpc.NewTCPConnection(url, &timeout)

	connection, err := baidurpc.NewTCPConnectionPool(url, &timeout, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer connection.Close()

	cc, cancel := context.WithDeadline(context.Background(), time.Now().Add(200*time.Second))

	for i := 0; i < *parell; i++ {
		go SendRpc(cc, connection)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	fmt.Println("Press Ctrl+C or send kill sinal to exit.")
	<-c
	cancel()
}

func SendRpc(c context.Context, connection baidurpc.Connection) {
	for {
		select {
		case <-c.Done():
			return
		default:
			time.Sleep(10 * time.Millisecond)
			doSimpleRPCInvoke(connection)
		}

	}

}

func doSimpleRPCInvoke(connection baidurpc.Connection) {

	// create client
	rpcClient, err := baidurpc.NewRpcCient(connection)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer rpcClient.Close()

	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	message := "say hello from xiemalin中文测试"
	dm := DataMessage{&message}

	rpcInvocation.SetParameterIn(&dm)
	rpcInvocation.LogId = proto.Int64(1)
	rpcInvocation.Attachment = []byte("hello world")

	parameterOut := DataMessage{}

	// 10000*time.Millisecond,
	response, err := rpcClient.SendRpcRequestWithTimeout(11*time.Millisecond, rpcInvocation, &parameterOut)
	if err != nil {
		fmt.Println(err)
		return
	}

	if response == nil {
		fmt.Println("Reponse is nil")
		return
	}

	fmt.Println("attachement", response.Attachment)

}

func Int(v int) *int {
	return &v
}

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
