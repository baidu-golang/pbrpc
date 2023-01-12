/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-07-24 16:54:14
 */
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
	"testing"

	baidurpc "github.com/baidu-golang/pbrpc"
	"google.golang.org/protobuf/proto"

	. "github.com/smartystreets/goconvey/convey"
)

// TestPropertySetAndGet
func TestPropertySetAndGet(t *testing.T) {
	Convey("TestPropertySetAndGet", t, func() {
		var serviceName string = "ThisAServiceName"
		var methodName string = "ThisAMethodName"
		var logId int64 = 1

		request := baidurpc.Request{
			ServiceName: serviceName,
			MethodName:  methodName,
			LogId:       logId,
		}

		So(serviceName, ShouldEqual, request.GetServiceName())
		So(methodName, ShouldEqual, request.GetMethodName())
		So(logId, ShouldEqual, request.GetLogId())

		data, err := proto.Marshal(&request)
		So(err, ShouldBeNil)

		request2 := new(baidurpc.Request)
		err = proto.Unmarshal(data, request2)
		So(err, ShouldBeNil)

		So(request.GetServiceName(), ShouldEqual, request2.GetServiceName())
	})

}
