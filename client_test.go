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
	"fmt"
	"log"

	baidurpc "github.com/baidu-golang/pbrpc"
	"github.com/golang/protobuf/proto"

	"strings"
	"testing"
	"time"
)

//bench mark test

func BenchmarkTestLocalConnectionPerformance(b *testing.B) {
	host := "localhost"
	port := 1039

	rpcServer := DoRpcServerStart(nil, port)
	defer rpcServer.Stop()

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 5
	// create client by simple connection
	connection, err := baidurpc.NewTCPConnection(url, &timeout)
	if err != nil {
		return
	}
	defer connection.Close()
	log.Println(b.N)
	for i := 0; i < b.N; i++ {
		doSimpleRPCInvoke(nil, connection)
	}

}

//bench mark test end

func TestSendRpcRequestSucc173300ess(t *testing.T) {

	host := "localhost"
	port := 1031

	rpcServer := DoRpcServerStart(t, port)
	defer rpcServer.Stop()

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 5
	// create client by simple connection
	connection, err := baidurpc.NewTCPConnection(url, &timeout)
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()
	doSimpleRPCInvoke(t, connection)

	// create client by connection pool
	connectionPool, err := baidurpc.NewDefaultTCPConnectionPool(url, &timeout)
	if err != nil {
		t.Error(err)
	}
	defer connectionPool.Close()
	doSimpleRPCInvoke(t, connectionPool)

}

func TestSendRpcRequestNoServiceNameExist(t *testing.T) {

	host := "localhost"
	port := 1035

	rpcServer := DoRpcServerStart(t, port)
	defer rpcServer.Stop()

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := time.Second * 5
	// create client by simple connection
	connection, err := baidurpc.NewTCPConnection(url, &timeout)
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()
	err = doSimpleRPCInvokeWithSignature(connection, "no", "no")
	if err == nil {
		t.Error("Service name not exist should return error.")
	}

}

func TestSendRpcRequestReconnectionSuccess(t *testing.T) {

	host := "localhost"
	port := 1030

	rpcServer := DoRpcServerStart(t, port)

	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)
	timeout := time.Second * 5
	// create client by connection pool
	connectionPool, err := baidurpc.NewDefaultTCPConnectionPool(url, &timeout)
	if err != nil {
		t.Error(err)
	}
	defer connectionPool.Close()
	doSimpleRPCInvoke(t, connectionPool)

	//close server
	rpcServer.Stop()

	DoRpcServerStart(t, port)
	doSimpleRPCInvoke(t, connectionPool)
}

func doSimpleRPCInvoke(t *testing.T, connection baidurpc.Connection) {

	serviceName := "echoService"
	methodName := "echo"

	doSimpleRPCInvokeWithSignature(connection, serviceName, methodName)
}

func doSimpleRPCInvokeWithSignature(connection baidurpc.Connection, serviceName, methodName string) error {

	// create client
	rpcClient, err := baidurpc.NewRpcCient(connection)
	if err != nil {
		return err
	}

	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)

	message := "say hello from xiemalin中文测试"
	dm := DataMessage{&message}

	rpcInvocation.SetParameterIn(&dm)
	rpcInvocation.LogId = proto.Int64(1)

	parameterOut := DataMessage{}

	response, err := rpcClient.SendRpcRequest(rpcInvocation, &parameterOut)
	if err != nil {
		return err
	}

	if response == nil {
		log.Println("Reponse is nil")
		return nil
	}

	expect := "hello " + message
	if !strings.EqualFold(expect, *parameterOut.Name) {
		return errors.New(fmt.Sprintf("expect message is '%s' but actual is '%s'", message, *parameterOut.Name))
	}
	return nil

}
