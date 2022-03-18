/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-04-26 18:18:59
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
package baidurpc_test

import (
	"testing"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
)

// TestHaClient test ha client
func TestHaClient(t *testing.T) {
	Convey("TestSingleTcpConnectionClient", t, func() {
		tcpServer := startRpcServer(0)
		defer stopRpcServer(tcpServer)

		host := "localhost"
		timeout := time.Second * 5
		errPort := 100
		urls := []baidurpc.URL{{Host: &host, Port: &errPort}, {Host: &host, Port: Int(PORT_1)}}

		connections, err := baidurpc.NewBatchTCPConnection(urls, timeout)
		So(err, ShouldBeNil)
		defer baidurpc.CloseBatchConnection(connections)

		haClient, err := baidurpc.NewHaRpcCient(connections)
		So(err, ShouldBeNil)
		So(haClient, ShouldNotBeNil)
		defer haClient.Close()

		Convey("Test send request EchoService!echo", func() {
			doHaSimpleRPCInvokeWithSignatureWithConvey(haClient, "EchoService", "echo", false, false, false, false)
		})
		Convey("Test send request async EchoService!echo", func() {
			doHaSimpleRPCInvokeWithSignatureWithConvey(haClient, "EchoService", "echo", false, false, true, false)
		})

	})
}

// doSimpleRPCInvokeWithSignatureWithConvey  send rpc request
func doHaSimpleRPCInvokeWithSignatureWithConvey(rpcClient *baidurpc.HaRpcClient, serviceName, methodName string, withAttachement, withCustomErr, async, timeout bool) {
	Convey("Test Client send rpc  request", func() {
		rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

		name := "马林"
		dm := DataMessage{Name: name}

		rpcInvocation.SetParameterIn(&dm)
		rpcInvocation.LogId = proto.Int64(1)

		if withAttachement {
			rpcInvocation.Attachment = []byte("This is attachment data")
		}

		parameterOut := DataMessage{}
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
