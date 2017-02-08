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
package pbrpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
)

var ERR_NEED_INIT = errors.New("[client-001]Session is not initialized, Please use NewRpcInvocation() to create instance.")
var ERR_RESPONSE_NIL = errors.New("[client-003]No response result, mybe net work break error.")
var LOG_SERVER_REPONSE_ERROR = "[client-002]Server response error. code=%d, msg='%s'"
var LOG_CLIENT_TIMECOUST_INFO = "[client-101]Server name '%s' method '%s' process cost '%.5g' seconds."

/*
RPC client invoke

*/
type RpcClient struct {
	Session Connection
}

type URL struct {
	Host *string
	Port *int
}

func (u *URL) SetHost(host *string) *URL {
	u.Host = host
	return u
}

func (u *URL) SetPort(port *int) *URL {
	u.Port = port
	return u
}

type RpcInvocation struct {
	ServiceName  *string
	MethodName   *string
	ParameterIn  *proto.Message
	Attachment   []byte
	LogId        *int64
	CompressType *int32
}

func NewRpcCient(connection Connection) (*RpcClient, error) {
	c := RpcClient{}
	c.Session = connection
	return &c, nil
}

func NewRpcInvocation(serviceName, methodName *string) *RpcInvocation {
	r := new(RpcInvocation)
	r.init(serviceName, methodName)

	return r
}

func (r *RpcInvocation) init(serviceName, methodName *string) {

	*r = RpcInvocation{}
	r.ServiceName = serviceName
	r.MethodName = methodName
	compressType := COMPRESS_NO
	r.CompressType = &compressType
	r.ParameterIn = nil
}

func (r *RpcInvocation) SetParameterIn(parameterIn proto.Message) {
	r.ParameterIn = &parameterIn
}

func (r *RpcInvocation) GetRequestRpcDataPackage() (*RpcDataPackage, error) {

	rpcDataPackage := new(RpcDataPackage)
	rpcDataPackage.ServiceName(*r.ServiceName)
	rpcDataPackage.MethodName(*r.MethodName)
	rpcDataPackage.MagicCode(MAGIC_CODE)
	rpcDataPackage.CompressType(*r.CompressType)
	rpcDataPackage.LogId(*r.LogId)

	if *r.ParameterIn != nil {
		data, err := proto.Marshal(*r.ParameterIn)
		if err != nil {
			return nil, err
		}
		rpcDataPackage.SetData(data)
	}

	return rpcDataPackage, nil
}

func (c *RpcClient) SendRpcRequest(rpcInvocation *RpcInvocation, responseMessage proto.Message) (*RpcDataPackage, error) {
	if c.Session == nil {
		return nil, ERR_NEED_INIT
	}

	now := time.Now().UnixNano()

	rpcDataPackage, err := rpcInvocation.GetRequestRpcDataPackage()
	if err != nil {
		return nil, err
	}

	rsp, err := c.Session.SendReceive(rpcDataPackage)
	if err != nil {
		return nil, err
	}

	r := rsp
	if r == nil {
		return nil, ERR_RESPONSE_NIL //to ingore this nil value
	}

	errorCode := r.GetMeta().GetResponse().GetErrorCode()
	if errorCode > 0 {
		errMsg := fmt.Sprintf(LOG_SERVER_REPONSE_ERROR,
			errorCode, r.GetMeta().GetResponse().GetErrorText())
		return nil, errors.New(errMsg)
	}

	response := r.GetData()
	if response != nil {
		err = proto.Unmarshal(response, responseMessage)
		if err != nil {
			return nil, err
		}
	}

	took := TimetookInSeconds(now)
	Infof(LOG_CLIENT_TIMECOUST_INFO, *rpcInvocation.ServiceName, *rpcInvocation.MethodName, took)

	return r, nil

}
