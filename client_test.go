// Go support for Protocol Buffers RPC which compatible with https://github.com/Baidu-ecom/Jprotobuf-rpc-socket
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
	"strings"
	"testing"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
)

const (
	AUTH_TOKEN = "SJIVNCQIN@#$@*sdjfsd"
)

type (
	StringMatchAuthService struct {
	}

	AddOneTraceService struct {
	}
)

// Authenticate
func (as *StringMatchAuthService) Authenticate(service, name string, authToken []byte) bool {
	if authToken == nil {
		return false
	}
	return strings.Compare(AUTH_TOKEN, string(authToken)) == 0
}

// Trace
func (as *AddOneTraceService) Trace(service, name string, traceInfo *baidurpc.TraceInfo) *baidurpc.TraceInfo {
	traceInfo.SpanId++
	traceInfo.TraceId++
	traceInfo.ParentSpanId++
	return traceInfo
}

// TestSingleTcpConnectionClient
func TestSingleTcpConnectionClient(t *testing.T) {
	Convey("TestSingleTcpConnectionClient", t, func() {
		tcpServer := startRpcServer(0)
		defer stopRpcServer(tcpServer)

		conn, client, err := createClient()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(client, ShouldNotBeNil)
		defer client.Close()
		defer conn.Close()

		testSendRpc("Client send rpc request", client, false, false, 0, false)
		testSendRpc("Client send rpc request(async)", client, true, false, 0, false)
	})
}

// TestSingleTcpConnectionClientWithAuthenticate
func TestSingleTcpConnectionClientWithAuthenticate(t *testing.T) {
	Convey("TestSingleTcpConnectionClientWithAuthenticate", t, func() {
		tcpServer := startRpcServer(0)
		tcpServer.SetAuthService(new(StringMatchAuthService))
		defer stopRpcServer(tcpServer)

		conn, client, err := createClient()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(client, ShouldNotBeNil)
		defer client.Close()
		defer conn.Close()

		testSendRpc("Client send rpc request", client, false, true, 0, false)
		testSendRpc("Client send rpc request(async)", client, true, true, 0, false)
	})
}

// TestSingleTcpConnectionClientWithChunk
func TestSingleTcpConnectionClientWithChunk(t *testing.T) {
	Convey("TestSingleTcpConnectionClientWithChunk", t, func() {
		tcpServer := startRpcServer(0)
		tcpServer.SetAuthService(new(StringMatchAuthService))
		defer stopRpcServer(tcpServer)

		conn, client, err := createClient()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(client, ShouldNotBeNil)
		defer client.Close()
		defer conn.Close()

		testSendRpc("Client send rpc request", client, false, true, 20, false)
		testSendRpc("Client send rpc request(async)", client, true, true, 20, false)
	})
}

// TestSingleTcpConnectionClientAndServerWithChunk
func TestSingleTcpConnectionClientAndServerWithChunk(t *testing.T) {
	Convey("TestSingleTcpConnectionClientAndServerWithChunk", t, func() {
		tcpServer := startRpcServer(20)
		tcpServer.SetAuthService(new(StringMatchAuthService))
		defer stopRpcServer(tcpServer)

		conn, client, err := createClient()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(client, ShouldNotBeNil)
		defer client.Close()
		defer conn.Close()

		testSendRpc("Client send rpc request", client, false, true, 20, false)
		testSendRpc("Client send rpc request(async)", client, true, true, 20, false)
	})
}

// TestSingleTcpConnectionClientWithBadChunkCase
func TestSingleTcpConnectionClientWithBadChunkCase(t *testing.T) {
	Convey("TestSingleTcpConnectionClientWithBadChunkCase", t, func() {
		tcpServer := startRpcServer(0)
		defer stopRpcServer(tcpServer)

		conn, client, err := createClient()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(client, ShouldNotBeNil)
		defer client.Close()
		defer conn.Close()

		serviceName := "EchoService"
		methodName := "echo"
		rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

		name := "(马林)(matthew)(XML)(jhunters)"
		dm := DataMessage{Name: name}

		rpcInvocation.SetParameterIn(&dm)
		rpcInvocation.LogId = proto.Int64(1)

		dataPackage, err := rpcInvocation.GetRequestRpcDataPackage()
		So(err, ShouldBeNil)

		dataPackage.ChuckInfo(10, 1) // bad chunk data package
		go func() {
			client.Session.SendReceive(dataPackage) // send bad chunk data package server will block unitl timeout
		}()
		time.Sleep(1 * time.Second)
		doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "echo", false, false, false, false, false, 0, false)
	})
}

// TestPooledTcpConnectionClient
func TestPooledTcpConnectionClient(t *testing.T) {
	Convey("TestPooledTcpConnectionClient", t, func() {
		tcpServer := startRpcServer(0)
		defer stopRpcServer(tcpServer)

		conn, client, err := createPooledConnectionClient()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(client, ShouldNotBeNil)
		defer client.Close()
		defer conn.Close()

		testSendRpc("Client send rpc request", client, false, true, 0, false)
		testSendRpc("Client send rpc request(async)", client, true, true, 0, false)
	})
}

// TestSingleTcpConnectionClientByAsync
func TestSingleTcpConnectionClientByAsync(t *testing.T) {
	Convey("TestSingleTcpConnectionClientByAsync", t, func() {
		tcpServer := startRpcServer(0)
		tcpServer.SetAuthService(new(StringMatchAuthService))
		defer stopRpcServer(tcpServer)

		conn, client, err := createClient()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(client, ShouldNotBeNil)
		defer client.Close()
		defer conn.Close()

		testSendRpc("Client send rpc request", client, false, true, 0, true)
		testSendRpc("Client send rpc request(async)", client, true, true, 0, true)
	})
}

func testSendRpc(testName string, client *baidurpc.RpcClient, timeout, auth bool, chunksize uint32, async bool) {
	Convey(testName, func() {
		Convey("Test send request EchoService!echo", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "echo", false, false, timeout, false, auth, chunksize, async)
		})
		Convey("Test send request EchoService!echoWithAttchement", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "echoWithAttchement", true, false, timeout, false, auth, chunksize, async)
		})
		Convey("Test send request EchoService!echoWithCustomizedError", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "echoWithCustomizedError", false, true, timeout, false, auth, chunksize, async)
		})
		Convey("Test send request EchoService!echoWithoutContext", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "echoWithoutContext", false, false, timeout, false, auth, chunksize, async)
		})
		Convey("Test send request EchoService!EchoSlowTest", func() {
			doSimpleRPCInvokeWithSignatureWithConvey(client, "EchoService", "EchoSlowTest", false, false, timeout, true, auth, chunksize, async)
		})
	})
}

// createRpcServer create rpc server by port and localhost
func createRpcServerWithChunkSize(port int, chunksize uint32) *baidurpc.TcpServer {
	serverMeta := baidurpc.ServerMeta{}
	serverMeta.Port = Int(port)
	serverMeta.ChunkSize = chunksize
	rpcServer := baidurpc.NewTpcServer(&serverMeta)
	return rpcServer
}

func startRpcServer(chunksize uint32) *baidurpc.TcpServer {
	return startRpcServerWithHttpMode(chunksize, false)
}

// startRpcServer start rpc server and register echo service as default rpc service
func startRpcServerWithHttpMode(chunksize uint32, httpMode bool) *baidurpc.TcpServer {

	rpcServer := createRpcServerWithChunkSize(PORT_1, chunksize)

	echoservice := new(EchoService)
	methodMapping := map[string]string{
		"Echo":                    "echo",
		"EchoWithAttchement":      "echoWithAttchement",
		"EchoWithCustomizedError": "echoWithCustomizedError",
		"EchoWithoutContext":      "echoWithoutContext",
	}
	rpcServer.RegisterNameWithMethodMapping("EchoService", echoservice, methodMapping)

	rpcServer.SetTraceService(new(AddOneTraceService))
	if httpMode {
		rpcServer.EnableHttp()
	}
	rpcServer.Start()

	return rpcServer
}

// createClient
func createClient() (baidurpc.Connection, *baidurpc.RpcClient, error) {

	host := "localhost"
	port := PORT_1

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 500
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
func doSimpleRPCInvokeWithSignatureWithConvey(rpcClient *baidurpc.RpcClient, serviceName, methodName string,
	withAttachement, withCustomErr, timeout, timeoutCheck, auth bool, chunkSize uint32, async bool) {
	Convey("Test Client send rpc  request", func() {
		rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

		name := "(马林)(matthew)(XML)(jhunters)"
		dm := DataMessage{Name: name}

		rpcInvocation.SetParameterIn(&dm)
		rpcInvocation.LogId = proto.Int64(1)
		rpcInvocation.ChunkSize = chunkSize
		rpcInvocation.TraceId = 10
		rpcInvocation.SpanId = 11
		rpcInvocation.ParentSpanId = 12
		rpcInvocation.RpcRequestMetaExt = map[string]string{"key1": "value1"}

		if withAttachement {
			rpcInvocation.Attachment = []byte("This is attachment data")
		}

		if auth {
			rpcInvocation.AuthenticateData = []byte(AUTH_TOKEN)
		}

		parameterOut := DataMessage{}
		var response *baidurpc.RpcDataPackage
		var err error
		if timeout {
			response, err = rpcClient.SendRpcRequestWithTimeout(1*time.Second, rpcInvocation, &parameterOut)
			if timeoutCheck {
				So(err, ShouldNotBeNil)
				return
			}
		} else {
			if async {
				ch := rpcClient.SendRpcRequestAsyc(rpcInvocation, &parameterOut)
				rpcResult := <-ch
				response = rpcResult.GetRpcDataPackage()
				err = rpcResult.GetErr()
			} else {
				response, err = rpcClient.SendRpcRequest(rpcInvocation, &parameterOut)
			}
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
			So(string(response.Attachment), ShouldEqual, "I am a attachement, This is attachment data")
		}

		So(response.GetTraceId(), ShouldEqual, rpcInvocation.TraceId+1)
		So(response.GetParentSpanId(), ShouldEqual, rpcInvocation.ParentSpanId+1)
		So(response.GetParentSpanId(), ShouldEqual, rpcInvocation.ParentSpanId+1)
		So(response.GetRpcRequestMetaExt()["key1"], ShouldEqual, "value1")
	})

}
