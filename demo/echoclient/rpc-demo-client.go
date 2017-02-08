package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	baidurpc "github.com/baidu-golang/baidurpc"

	"github.com/golang/protobuf/proto"
	pool "github.com/jolestar/go-commons-pool"
)

var host = flag.String("host", "localhost", "If non-empty, connect to target host")
var port = flag.Int("port", 8122, "If non-empty, connect to target host")

var countsPerThread = flag.Int("countsPerThread", 100, "If non-empty, counts per thread to execute default is 100")
var parell = flag.Int("parell", 5, "If non-empty, concurrent threads to run")

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

func main() {

	url := baidurpc.URL{}
	url.SetHost(host).SetPort(port)

	timeout := time.Second * 5
	// create client by connection pool
	config := pool.NewDefaultPoolConfig()
	config.MaxTotal = *parell
	config.MaxIdle = *parell
	connection, err := baidurpc.NewTCPConnectionPool(url, &timeout, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer connection.Close()

	ch := make(chan int, *parell)
	now := time.Now().Unix()
	for y := 0; y < *parell; y++ {

		go SendRpc(connection, ch, y)
	}

	var count = 0
	for {
		_ = <-ch
		fmt.Println("Receive chan back")
		count++
		if count > (*parell - 1) {
			break
		}
	}

	after := time.Now().Unix()

	timecost := int(after - now)
	if timecost == 0 {
		timecost = 1
	}
	cpt := *countsPerThread
	pl := *parell
	fmt.Println("Performance:", cpt*pl/timecost, " requests in secend.")

	fmt.Println(after - now)
}

func SendRpc(connection baidurpc.Connection, ch chan int, y int) {
	fmt.Println("start go", y)
	for i := 0; i < *countsPerThread; i++ {
		doSimpleRPCInvoke(connection, i, y)

	}
	fmt.Println("end go", y)
	ch <- y
}

func doSimpleRPCInvoke(connection baidurpc.Connection, x, y int) {

	// create client
	rpcClient, err := baidurpc.NewRpcCient(connection)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	//defer rpcClient.Close()

	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	// 指定压缩算法
	rpcInvocation.CompressType = proto.Int32(baidurpc.COMPRESS_GZIP)

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
