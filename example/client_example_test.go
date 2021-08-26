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
	"fmt"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"
	pool "github.com/jolestar/go-commons-pool/v2"
	"google.golang.org/protobuf/proto"
)

// ExampleTcpClient
func ExampleTcpClient() {

	host := "localhost"
	port := 1031

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 5
	// create client by connection pool
	connection, err := baidurpc.NewTCPConnection(url, &timeout)
	if err != nil {
		return
	}
	defer connection.Close()

	// create client
	rpcClient, err := baidurpc.NewRpcCient(connection)
	if err != nil {
		return
	}
	defer rpcClient.Close()

	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	message := "say hello from xiemalin中文测试"
	dm := DataMessage{Name: message}

	rpcInvocation.SetParameterIn(&dm)
	rpcInvocation.LogId = proto.Int64(1)
	rpcInvocation.Attachment = []byte("hello world")

	parameterOut := DataMessage{}

	response, err := rpcClient.SendRpcRequest(rpcInvocation, &parameterOut)
	if err != nil {
		return
	}

	if response == nil {
		fmt.Println("Reponse is nil")
		return
	}

	fmt.Println("attachement", response.Attachment)

}

// ExampleTcpClient
func ExamplePooledTcpClient() {

	host := "localhost"
	port := 1031

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 5

	parell := 2

	// create client by connection pool
	config := pool.NewDefaultPoolConfig()
	config.MaxTotal = parell
	config.MaxIdle = parell

	// create client by connection pool
	connection, err := baidurpc.NewTCPConnectionPool(url, &timeout, config)
	if err != nil {
		return
	}
	defer connection.Close()

	// create client
	rpcClient, err := baidurpc.NewRpcCient(connection)
	if err != nil {
		return
	}
	defer rpcClient.Close()

	serviceName := "echoService"
	methodName := "echo"
	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	message := "say hello from xiemalin中文测试"
	dm := DataMessage{Name: message}

	rpcInvocation.SetParameterIn(&dm)
	rpcInvocation.LogId = proto.Int64(1)
	rpcInvocation.Attachment = []byte("hello world")

	parameterOut := DataMessage{}

	response, err := rpcClient.SendRpcRequest(rpcInvocation, &parameterOut)
	if err != nil {
		return
	}

	if response == nil {
		fmt.Println("Reponse is nil")
		return
	}

	fmt.Println("attachement", response.Attachment)

}
