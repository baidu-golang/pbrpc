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
	"testing"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"
	"github.com/golang/protobuf/proto"
	. "github.com/smartystreets/goconvey/convey"
)

// TestSingleTcpConnectionClient
func TestSingleTcpConnectionClient(t *testing.T) {
	Convey("TestSingleTcpConnectionClient", t, func() {
		tcpServer := startRpcServer()
		defer tcpServer.Stop()

		conn, client, err := createClient()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(client, ShouldNotBeNil)
		defer client.Close()
		defer conn.Close()

		testSendRpc("Client send rpc request", client, true)
		testSendRpc("Client send rpc request(async)", client, false)
	})
}

// TestPooledTcpConnectionClient
func TestPooledTcpConnectionClient(t *testing.T) {
	Convey("TestSingleTcpConnectionClient", t, func() {
		tcpServer := startRpcServer()
		defer tcpServer.Stop()

		conn, client, err := createPooledConnectionClient()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(client, ShouldNotBeNil)
		defer client.Close()
		defer conn.Close()

		testSendRpc("Client send rpc request", client, true)
		testSendRpc("Client send rpc request(async)", client, false)
	})
}

func testSendRpc(testName string, client *baidurpc.RpcClient, async bool) {
	Convey(testName, func() {
		Convey("Test send request EchoService!echo", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "echo", false, false, async, false)
		})
		Convey("Test send request EchoService!echoWithAttchement", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "echoWithAttchement", true, false, async, false)
		})
		Convey("Test send request EchoService!echoWithCustomizedError", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "echoWithCustomizedError", false, true, async, false)
		})
		Convey("Test send request EchoService!echoWithoutContext", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "echoWithoutContext", false, false, async, false)
		})
		Convey("Test send request EchoService!EchoSlowTest", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "EchoSlowTest", false, false, async, true)
		})
	})
}

// startRpcServer start rpc server and register echo service as default rpc service
func startRpcServer() *baidurpc.TcpServer {

	rpcServer := createRpcServer(PORT_1)

	echoservice := new(baidurpc.EchoService)
	methodMapping := map[string]string{
		"Echo":                    "echo",
		"EchoWithAttchement":      "echoWithAttchement",
		"EchoWithCustomizedError": "echoWithCustomizedError",
		"EchoWithoutContext":      "echoWithoutContext",
	}
	rpcServer.RegisterNameWithMethodMapping("EchoService", echoservice, methodMapping)

	rpcServer.Start()

	return rpcServer
}

// createClient
func createClient() (baidurpc.Connection, *baidurpc.RpcClient, error) {

	host := "localhost"
	port := PORT_1

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 5
	// create client by simple connection
	connection, err := baidurpc.NewTCPConnection(url, &timeout)
	if err != nil {
		return nil, nil, err
	}
	rpcClient, err := baidurpc.NewRpcCient(connection)
	if err != nil {
		return nil, nil, err
	}
	return connection, rpcClient, nil
}

// createClient
func createPooledConnectionClient() (baidurpc.Connection, *baidurpc.RpcClient, error) {

	host := "localhost"
	port := PORT_1

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 5
	// create client by simple connection
	connection, err := baidurpc.NewTCPConnectionPool(url, &timeout, nil)
	if err != nil {
		return nil, nil, err
	}
	rpcClient, err := baidurpc.NewRpcCient(connection)
	if err != nil {
		return nil, nil, err
	}
	return connection, rpcClient, nil
}

// doSimpleRPCInvokeWithSignatureWithConvey  send rpc request
func doSimpleRPCInvokeWithSignatureWithConvey(rpcClient *baidurpc.RpcClient, serviceName, methodName string, withAttachement, withCustomErr, async, timeout bool) {
	Convey("Test Client send rpc  request", func() {
		rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

		name := "马林"
		dm := baidurpc.EchoMessage{name}

		rpcInvocation.SetParameterIn(&dm)
		rpcInvocation.LogId = proto.Int64(1)

		if withAttachement {
			rpcInvocation.Attachment = []byte("This is attachment data")
		}

		parameterOut := baidurpc.EchoMessage{}
		var response *baidurpc.RpcDataPackage
		var err error
		if async {
			response, err = rpcClient.SendRpcRequestWithTimeout(1*time.Second, rpcInvocation, &parameterOut)
			if timeout {
				So(err, ShouldNotBeNil)
				return
			}
		} else {
			response, err = rpcClient.SendRpcRequest(rpcInvocation, &parameterOut)
		}
		if withCustomErr {
			So(err, ShouldNotBeNil)
			return
		} else {
			So(err, ShouldBeNil)
		}
		So(response, ShouldNotBeNil)
		expect := "hello " + name
		So(expect, ShouldEqual, parameterOut.Name)

		if withAttachement {
			So(string(response.Attachment), ShouldEqual, "I am a attachementThis is attachment data")
		}

	})

}
