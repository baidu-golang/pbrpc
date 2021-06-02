package main

import (
	"flag"
	"os"

	baidurpc "github.com/baidu-golang/pbrpc"
	"github.com/golang/protobuf/proto"
)

var port = flag.Int("port", 8122, "If non-empty, port this server to listen")

func init() {

	if !flag.Parsed() {
		flag.Parse()
	}
}

func main() {

	serverMeta := baidurpc.ServerMeta{}
	serverMeta.Host = nil
	serverMeta.Port = Int(*port)
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	echoService := new(EchoService)

	rpcServer.RegisterName("echoService", echoService)

	mapping := make(map[string]string)
	mapping["Echo"] = "echo"
	rpcServer.RegisterNameWithMethodMapping("echoService", echoService, mapping)

	err := rpcServer.StartAndBlock()

	if err != nil {
		baidurpc.Error(err)
		os.Exit(-1)
	}
}

type EchoService struct {
}

func Int(v int) *int {
	return &v
}

func (rpc *EchoService) Echo(in *DataMessage) *DataMessage {
	var ret = "hello "

	if len(*in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + *in.Name
	}
	dm := DataMessage{}
	dm.Name = proto.String(ret)
	return &dm
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
