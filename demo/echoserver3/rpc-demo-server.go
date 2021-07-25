/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-06-04 14:25:31
 */
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"
	"github.com/baidu-golang/pbrpc/nettool"
	"github.com/golang/protobuf/proto"
)

var port = flag.Int("port", 8122, "If non-empty, port this server to listen")

func init() {

	if !flag.Parsed() {
		flag.Parse()
	}
}

type info struct {
	sTime time.Time
}

func main() {

	serverMeta := baidurpc.ServerMeta{}
	serverMeta.Host = nil
	serverMeta.Port = Int(*port)
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	echoService := new(EchoService)

	// rpcServer.RegisterName("echoService", echoService)

	mapping := make(map[string]string)
	mapping["Echo"] = "echo"
	rpcServer.RegisterNameWithMethodMapping("echoService", echoService, mapping)

	addr := ":" + strconv.Itoa(*port)
	var headsize uint8 = 9
	selector, err := nettool.NewCustomListenerSelectorByAddr(addr, headsize, nettool.StartWith_Mode)
	if err != nil {
		fmt.Println(err)
		return
	}

	rpcServerListener, err := selector.RegisterListener(baidurpc.MAGIC_CODE) //"PRPC"
	if err != nil {
		fmt.Println(err)
		return
	}
	httpServerListener := selector.RegisterDefaultListener()

	hsv := HttpStatusView{&info{sTime: time.Now()}}
	srv := &http.Server{
		Handler: &hsv,
	}
	go func() {
		// service connections
		if err := srv.Serve(httpServerListener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go selector.Serve()

	rpcServer.StartServer(rpcServerListener)

	if err != nil {
		baidurpc.Error(err)
		os.Exit(-1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	fmt.Println("Press Ctrl+C or send kill sinal to exit.")
	<-c
}

type HttpStatusView struct {
	server *info
}

func (hsv *HttpStatusView) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	v, _ := json.Marshal(hsv.server.sTime)
	w.Write(v)
}

type EchoService struct {
}

func Int(v int) *int {
	return &v
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) Echo(c context.Context, in *DataMessage) (*DataMessage, context.Context) {
	var ret = "hello "

	attachement := baidurpc.Attachement(c)
	fmt.Println("attachement", attachement)

	if len(*in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + *in.Name
	}
	dm := DataMessage{}
	dm.Name = proto.String(ret)

	// bind attachment
	cc := baidurpc.BindAttachement(context.Background(), []byte("hello"))
	// bind with err
	// cc = baidurpc.BindError(cc, errors.New("manule error"))
	return &dm, cc
}

// EchoWithoutContext
func (rpc *EchoService) EchoWithoutContext(c context.Context, in *DataMessage) *DataMessage {
	dm, _ := rpc.Echo(c, in)
	return dm
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
