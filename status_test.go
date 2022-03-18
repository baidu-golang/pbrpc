/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-04 10:43:27
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
 * @Date: 2021-08-04 10:43:27
 */
package baidurpc_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
)

const (
	Default_Timeout = 10 * time.Second
)

// TestRpcStatus
func TestRpcStatus(t *testing.T) {
	Convey("TestRpcStatus", t, func() {
		tcpServer := startRpcServer(0)
		defer stopRpcServer(tcpServer)

		serviceName := baidurpc.RPC_STATUS_SERVICENAME
		methodName := "Status"

		parameterOut := &baidurpc.RPCStatus{}
		err := sendRpc("localhost", PORT_1, serviceName, methodName, nil, parameterOut)
		So(err, ShouldBeNil)
		So(parameterOut.Port, ShouldEqual, PORT_1)
		echoservice := new(EchoService)
		tt := reflect.TypeOf(echoservice)

		So(len(parameterOut.Methods), ShouldEqual, tt.NumMethod()+2)

	})
}

// TestRpcRequestStatus
func TestRpcRequestStatus(t *testing.T) {
	Convey("TestRpcStatus", t, func() {
		tcpServer := startRpcServer(0)
		defer stopRpcServer(tcpServer)

		serviceName := baidurpc.RPC_STATUS_SERVICENAME
		methodName := "QpsDataStatus"

		parameterOut := &baidurpc.QpsData{}
		parameterIn := &baidurpc.RPCMethod{Service: baidurpc.RPC_STATUS_SERVICENAME, Method: methodName}
		err := sendRpc("localhost", PORT_1, serviceName, methodName, parameterIn, parameterOut)
		So(err, ShouldBeNil)

	})
}

func sendRpc(host string, port int, serviceName, methodName string, parameterIn, parameterOut proto.Message) error {
	url := baidurpc.URL{}
	url.SetHost(&host).SetPort(&port)

	timeout := Default_Timeout
	connection, err := baidurpc.NewTCPConnection(url, &timeout)
	if err != nil {
		return err
	}
	defer connection.Close()

	// create client
	rpcClient, err := baidurpc.NewRpcCient(connection)
	if err != nil {
		return err
	}

	rpcInvocation := baidurpc.NewRpcInvocation(&serviceName, &methodName)
	if parameterIn != nil {
		rpcInvocation.SetParameterIn(parameterIn)
	}

	rpcDataPackage, err := rpcClient.SendRpcRequest(rpcInvocation, parameterOut)
	if int(rpcDataPackage.Meta.GetResponse().ErrorCode) == baidurpc.ST_SERVICE_NOTFOUND {
		return fmt.Errorf("remote server not support this feature,  please upgrade version")
	}

	return err
}
