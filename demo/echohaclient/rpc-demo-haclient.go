package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"

	"github.com/golang/protobuf/proto"
)

var host = flag.String("host", "localhost", "If non-empty, connect to target host")
var port = flag.Int("port", 8122, "If non-empty, use this port")

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

func main() {

	timeout := time.Second * 5
	errPort := 100
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

	now := time.Now().Unix()
	for i := 0; i < 20; i++ {
		doSimpleRPCInvoke(haClient)
	}
	after := time.Now().Unix()

	timecost := int(after - now)
	if timecost == 0 {
		timecost = 1
	}

	fmt.Println(after - now)
}

func doSimpleRPCInvoke(rpcClient *baidurpc.HaRpcClient) {

	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	message := "say hello from xiemalin中文测试"
	dm := DataMessage{&message}

	rpcInvocation.SetParameterIn(&dm)
	rpcInvocation.LogId = proto.Int64(1)
	rpcInvocation.Attachment = []byte("hello world")

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
