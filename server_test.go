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
package baidurpc_test

import (
	"errors"
	"testing"

	baidurpc "github.com/baidu-golang/baidurpc"
	"github.com/golang/protobuf/proto"
)

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

		m, ok := msg.(*DataMessage)
		if !ok {
			errStr := "message type is not type of 'DataMessage'"
			return nil, nil, errors.New(errStr)
		}
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

func TestRpcServerStart(t *testing.T) {

	if true {
		return
	}
	DoRpcServerStartAndBlock(t, 1032)

}

func DoRpcServerStartAndBlock(t *testing.T, port int) {

	rpcServer := createRpcServer(port)

	err := rpcServer.StartAndBlock()
	if err != nil {
		t.Error(err)
	}
}

func DoRpcServerStart(t *testing.T, port int) *baidurpc.TcpServer {

	rpcServer := createRpcServer(port)

	err := rpcServer.Start()
	if err != nil && t != nil {
		t.Error(err)
	}
	return rpcServer
}

func createRpcServer(port int) *baidurpc.TcpServer {
	serverMeta := baidurpc.ServerMeta{}
	serverMeta.Host = proto.String("localhost")
	serverMeta.Port = Int(port)
	rpcServer := baidurpc.NewTpcServer(&serverMeta)

	ss := NewSimpleService("echoService", "echo")

	rpcServer.Register(ss)

	return rpcServer
}

func Int(v int) *int {
	p := new(int)
	*p = int(v)
	return p
}
