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

	baidurpc "github.com/baidu-golang/pbrpc"
	"github.com/golang/protobuf/proto"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	PORT_1 = 1031
	PORT_2 = 1032
	PORT_3 = 1033
)

// createRpcServer create rpc server by port and localhost
func createRpcServer(port int) *baidurpc.TcpServer {
	serverMeta := baidurpc.ServerMeta{}
	serverMeta.Port = Int(port)
	rpcServer := baidurpc.NewTpcServer(&serverMeta)
	return rpcServer
}

type CustomAuthService struct {
}

// Authenticate
func (as *CustomAuthService) Authenticate(service, name string, authToken []byte) bool {
	return true
}

// TestServerWithAuthenticate
func TestServerWithAuthenticate(t *testing.T) {
	Convey("TestServerWithoutPublishMethods", t, func() {

		rpcServer := createRpcServer(PORT_2)
		authservice := new(CustomAuthService)
		rpcServer.SetAuthService(authservice)
		err := rpcServer.Start()
		So(err, ShouldBeNil)
		rpcServer.Stop()
		So(err, ShouldBeNil)
	})
}

// TestServerWithoutPublishMethods
func TestServerWithoutPublishMethods(t *testing.T) {
	Convey("TestServerWithoutPublishMethods", t, func() {

		rpcServer := createRpcServer(PORT_2)
		err := rpcServer.Start()
		So(err, ShouldBeNil)
		rpcServer.Stop()
		So(err, ShouldBeNil)
	})
}

// TestServerWithPublishMethods
func TestServerWithPublishMethods(t *testing.T) {
	Convey("TestServerWithoutPublishMethods", t, func() {

		Convey("publish method with RegisterName", func() {
			rpcServer := createRpcServer(PORT_2)

			echoservice := new(baidurpc.EchoService)
			rpcServer.RegisterName("EchoService", echoservice)

			err := rpcServer.Start()
			So(err, ShouldBeNil)
			rpcServer.Stop()
			So(err, ShouldBeNil)
		})

		Convey("publish method with RegisterNameWithMethodMapping", func() {
			rpcServer := createRpcServer(PORT_2)

			echoservice := new(baidurpc.EchoService)
			methodMapping := map[string]string{
				"Echo":                    "echo",
				"EchoWithAttchement":      "echoWithAttchement",
				"EchoWithCustomizedError": "echoWithCustomizedError",
				"EchoWithoutContext":      "echoWithoutContext",
			}
			rpcServer.RegisterNameWithMethodMapping("EchoService", echoservice, methodMapping)

			err := rpcServer.Start()
			So(err, ShouldBeNil)
			rpcServer.Stop()
			So(err, ShouldBeNil)
		})

	})
}

// SimpleService
type SimpleService struct {
	serviceName string
	methodName  string
}

// NewSimpleService create RPC service
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

// TestServerWithOldRegisterWay
func TestServerWithOldRegisterWay(t *testing.T) {

	Convey("TestServerWithOldRegisterWay", t, func() {
		rpcServer := createRpcServer(PORT_3)
		So(rpcServer, ShouldNotBeNil)

		service := NewSimpleService("OldService", "OldMethod")
		rpcServer.Register(service)

		err := rpcServer.Start()
		So(err, ShouldBeNil)
		rpcServer.Stop()
		So(err, ShouldBeNil)
	})

}

// Int convert to pointer type of int
func Int(v int) *int {
	p := new(int)
	*p = int(v)
	return p
}
