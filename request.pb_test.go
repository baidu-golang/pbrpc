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
	baidurpc "github.com/baidu-golang/pbrpc"
	"github.com/golang/protobuf/proto"
	"strings"
	"testing"
)

func TestPropertySetAndGet(t *testing.T) {
    
    if true {return}

	var serviceName string = "ThisAServiceName"
	var methodName string = "ThisAMethodName"
	var logId int64 = 1

	request := baidurpc.Request{
		ServiceName: &serviceName,
		MethodName:  &methodName,
		LogId:       &logId,
	}

	if !strings.EqualFold(serviceName, request.GetServiceName()) {
		t.Errorf("set ServiceName value is '%s', but get value is '%s' ", serviceName, request.GetServiceName())
	}

	if !strings.EqualFold(methodName, request.GetMethodName()) {
		t.Errorf("set methodName value is '%s', but get value is '%s' ", methodName, request.GetMethodName())
	}

	if logId != request.GetLogId() {
		t.Errorf("set logId value is '%d', but get value is '%d' ", logId, request.GetLogId())
	}

	data, err := proto.Marshal(&request)
	if err != nil {
		t.Errorf("marshaling error: %s", err.Error())
	}
	
	request2 := new(baidurpc.Request)
	err = proto.Unmarshal(data, request2)
	if (err != nil) {
		t.Errorf("marshaling error: %s", err.Error())
	}
	
	if !strings.EqualFold(request.GetServiceName(), request2.GetServiceName()) {
		t.Errorf("set ServiceName value is '%s', but get value is '%s' ", request.GetServiceName(), request2.GetServiceName())
	}
}
