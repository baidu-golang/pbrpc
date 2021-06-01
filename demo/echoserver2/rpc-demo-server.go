package main

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"strings"

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

	rpcServer := createRpcServer(*port)

	err := rpcServer.StartAndBlock()

	if err != nil {
		baidurpc.Error(err)
		os.Exit(-1)
	}

}

func createRpcServer(port int) *baidurpc.TcpServer {
	serverMeta := baidurpc.ServerMeta{}
	serverMeta.Host = nil
	serverMeta.Port = Int(port)
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	callback := func(msg proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error) {
		var ret = "hello "

		if msg != nil {
			t := reflect.TypeOf(msg)

			if !strings.Contains(t.String(), "DataMessage") {
				errStr := "message type is not type of 'DataMessage'" + t.String()
				return nil, nil, errors.New(errStr)
			}

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

	rpcServer.RegisterRpc("echoService", "echo", callback, &DataMessage{})

	return rpcServer
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
