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
 * @Date: 2021-08-03 15:52:08
 */
package baidurpc_test

import (
	"testing"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"
	"github.com/golang/protobuf/proto"
)

// BenchmarkTestLocalConnectionPerformance bench mark test
func BenchmarkTestLocalConnectionPerformance(b *testing.B) {
	host := "localhost"
	port := PORT_1

	tcpServer := startRpcServer()
	defer tcpServer.Stop()

	conn, client, _ := createClient()
	defer client.Close()
	defer conn.Close()

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 5
	// create client by simple connection
	connection, err := baidurpc.NewTCPConnection(url, &timeout)
	if err != nil {
		return
	}
	defer connection.Close()
	for i := 0; i < b.N; i++ {
		doSimpleRPCInvokeWithSignature(client, "EchoService", "echo")
	}
}

// doSimpleRPCInvokeWithSignature  send rpc request
func doSimpleRPCInvokeWithSignature(rpcClient *baidurpc.RpcClient, serviceName, methodName string) {
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	name := "马林"
	dm := EchoMessage{name}

	rpcInvocation.SetParameterIn(&dm)
	rpcInvocation.LogId = proto.Int64(1)

	parameterOut := EchoMessage{}

	rpcClient.SendRpcRequest(rpcInvocation, &parameterOut)
}
