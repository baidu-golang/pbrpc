/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-03 19:30:12
 */
// Go support for Protocol Buffers RPC which compatiable with https://github.com/Baidu-ecom/Jprotobuf-rpc-socket
//
// Copyright 2002-2007 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-03 19:12:14
 */
package example

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"reflect"
	"strings"

	baidurpc "github.com/baidu-golang/pbrpc"
	"google.golang.org/protobuf/proto"
)

// ExampleRpcServer
func ExampleRpcServer() {
	port := 1031

	serverMeta := baidurpc.ServerMeta{}
	serverMeta.QPSExpireInSecs = 600 // set qps monitor expire time
	serverMeta.Port = &port
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	echoService := new(EchoService)

	mapping := make(map[string]string)
	mapping["Echo"] = "echo"
	rpcServer.RegisterNameWithMethodMapping("echoService", echoService, mapping)

	rpcServer.Start()
	defer rpcServer.Stop()

}

func ExampleRpcServerWithHttp() {
	port := 1031

	serverMeta := baidurpc.ServerMeta{}
	serverMeta.QPSExpireInSecs = 600 // set qps monitor expire time
	serverMeta.Port = &port
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	echoService := new(EchoService)

	mapping := make(map[string]string)
	mapping["Echo"] = "echo"
	rpcServer.RegisterNameWithMethodMapping("echoService", echoService, mapping)

	// start http rpc mode
	rpcServer.EnableHttp()
	rpcServer.Start()
	defer rpcServer.Stop()

}

// ExampleRpcServerWithListener
func ExampleRpcServerWithListener() {
	serverMeta := baidurpc.ServerMeta{}
	serverMeta.QPSExpireInSecs = 600 // set qps monitor expire time
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	echoService := new(EchoService)

	mapping := make(map[string]string)
	mapping["Echo"] = "echo"
	rpcServer.RegisterNameWithMethodMapping("echoService", echoService, mapping)

	addr := ":1031"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	rpcServer.StartServer(listener)
	defer rpcServer.Stop()

}

// ExampleRpcServerWithListener
func ExampleRpcServerRegisterWithCallback() {
	serverMeta := baidurpc.ServerMeta{}
	serverMeta.QPSExpireInSecs = 600 // set qps monitor expire time
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	callback := func(msg proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error) {
		var ret = "hello "

		if msg != nil {
			t := reflect.TypeOf(msg)

			if !strings.Contains(t.String(), "DataMessage") {
				errStr := "message type is not type of 'DataMessage'" + t.String()
				return nil, nil, errors.New(errStr)
			}

			var name string

			m := msg.(*DataMessage)
			name = m.Name

			if len(name) == 0 {
				ret = ret + "veryone"
			} else {
				ret = ret + name
			}
		}
		dm := DataMessage{}
		dm.Name = ret
		return &dm, []byte{1, 5, 9}, nil
	}

	rpcServer.RegisterRpc("echoService", "echo", callback, &DataMessage{})

	rpcServer.Start()
	defer rpcServer.Stop()

}

type EchoService struct {
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) Echo(c context.Context, in *DataMessage) (*DataMessage, context.Context) {
	var ret = "hello "

	attachement := baidurpc.Attachement(c)
	fmt.Println("attachement", attachement)

	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}
	dm := DataMessage{}
	dm.Name = ret

	// bind attachment
	cc := baidurpc.BindAttachement(context.Background(), []byte("hello"))
	// bind with err
	// cc = baidurpc.BindError(cc, errors.New("manule error"))
	return &dm, cc
}

type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

func DoRPCServer() {
	fmt.Println("server starting...")
	rpc.RegisterName("HelloService", new(HelloService)) // 注册RPC， RPC名称为 HelloService

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("Accept error:", err)
	}

	rpc.ServeConn(conn) // 绑定tcp服务链接

	fmt.Println("server started.")
}
