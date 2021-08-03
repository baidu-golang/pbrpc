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
package baidurpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jhunters/timewheel"
)

var (
	defaultTimewheelInterval = 10 * time.Millisecond
	defaultTimewheelSlot     = 300

	ERR_NEED_INIT             = errors.New("[client-001]Session is not initialized, Please use NewRpcInvocation() to create instance")
	ERR_RESPONSE_NIL          = errors.New("[client-003]No response result, mybe net work break error")
	LOG_SERVER_REPONSE_ERROR  = "[client-002]Server response error. code=%d, msg='%s'"
	LOG_CLIENT_TIMECOUST_INFO = "[client-101]Server name '%s' method '%s' process cost '%.5g' seconds"
)

const (
	ST_READ_TIMEOUT = 62
)

/*
RPC client invoke
*/
type RpcClient struct {
	Session Connection
	tw      *timewheel.TimeWheel
}

type URL struct {
	Host *string
	Port *int
}

// SetHost set host name
func (u *URL) SetHost(host *string) *URL {
	u.Host = host
	return u
}

// SetPort set port
func (u *URL) SetPort(port *int) *URL {
	u.Port = port
	return u
}

// RpcInvocation define rpc invocation
type RpcInvocation struct {
	ServiceName  *string
	MethodName   *string
	ParameterIn  *proto.Message
	Attachment   []byte
	LogId        *int64
	CompressType *int32
}

// NewRpcCient new rpc client
func NewRpcCient(connection Connection) (*RpcClient, error) {
	return NewRpcCientWithTimeWheelSetting(connection, defaultTimewheelInterval, uint16(defaultTimewheelSlot))
}

// NewRpcCientWithTimeWheelSetting new rpc client with set timewheel settings
func NewRpcCientWithTimeWheelSetting(connection Connection, timewheelInterval time.Duration, timewheelSlot uint16) (*RpcClient, error) {
	c := RpcClient{}
	c.Session = connection
	c.tw, _ = timewheel.New(timewheelInterval, timewheelSlot)
	c.tw.Start()
	return &c, nil
}

// NewRpcInvocation create RpcInvocation with service name and method name
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

// SetParameterIn
func (r *RpcInvocation) SetParameterIn(parameterIn proto.Message) {
	r.ParameterIn = &parameterIn
}

// GetRequestRpcDataPackage
func (r *RpcInvocation) GetRequestRpcDataPackage() (*RpcDataPackage, error) {

	rpcDataPackage := new(RpcDataPackage)
	rpcDataPackage.ServiceName(*r.ServiceName)
	rpcDataPackage.MethodName(*r.MethodName)
	rpcDataPackage.MagicCode(MAGIC_CODE)
	if r.CompressType != nil {
		rpcDataPackage.CompressType(*r.CompressType)
	}
	if r.LogId != nil {
		rpcDataPackage.LogId(*r.LogId)
	}

	rpcDataPackage.SetAttachment(r.Attachment)

	if r.ParameterIn != nil {
		data, err := proto.Marshal(*r.ParameterIn)
		if err != nil {
			return nil, err
		}
		rpcDataPackage.SetData(data)
	}

	return rpcDataPackage, nil
}

// define client methods
// Close close client with time wheel
func (c *RpcClient) Close() {
	if c.tw != nil {
		c.tw.Stop()
	}
}

// asyncRequest
func (c *RpcClient) asyncRequest(timeout time.Duration, request *RpcDataPackage, ch chan<- *RpcDataPackage) {
	// create a task bind with key, data and  time out call back function.
	t := &timewheel.Task{
		Data: nil, // business data
		TimeoutCallback: func(task timewheel.Task) { // call back function on time out
			// process someting after time out happened.
			errorcode := int32(ST_READ_TIMEOUT)
			request.ErrorCode(errorcode)
			errormsg := fmt.Sprintf("request time out of %v", task.Delay())
			request.ErrorText(errormsg)
			ch <- request
		}}

	// add task and return unique task id
	taskid, err := c.tw.AddTask(timeout, *t) // add delay task
	if err != nil {
		errorcode := int32(ST_ERROR)
		request.ErrorCode(errorcode)
		errormsg := err.Error()
		request.ErrorText(errormsg)

		ch <- request
		return
	}

	defer func() {
		c.tw.RemoveTask(taskid)
		if e := recover(); e != nil {
			Warningf("asyncRequest failed with error %v", e)
		}
	}()

	rsp, err := c.Session.SendReceive(request)
	if err != nil {
		errorcode := int32(ST_ERROR)
		request.ErrorCode(errorcode)
		errormsg := err.Error()
		request.ErrorText(errormsg)

		ch <- request
		return
	}

	ch <- rsp
}

// SendRpcRequest send rpc request to remote server
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
		errorcode := int32(ST_ERROR)
		rpcDataPackage.ErrorCode(errorcode)
		errormsg := err.Error()
		rpcDataPackage.ErrorText(errormsg)
		return rpcDataPackage, err
	}

	r := rsp
	if r == nil {
		return nil, ERR_RESPONSE_NIL //to ingore this nil value
	}

	errorCode := r.GetMeta().GetResponse().GetErrorCode()
	if errorCode > 0 {
		errMsg := fmt.Sprintf(LOG_SERVER_REPONSE_ERROR,
			errorCode, r.GetMeta().GetResponse().GetErrorText())
		return r, errors.New(errMsg)
	}

	response := r.GetData()
	if response != nil {
		err = proto.Unmarshal(response, responseMessage)
		if err != nil {
			return r, err
		}
	}

	took := TimetookInSeconds(now)
	Infof(LOG_CLIENT_TIMECOUST_INFO, *rpcInvocation.ServiceName, *rpcInvocation.MethodName, took)

	return r, nil

}

// SendRpcRequest send rpc request to remote server
func (c *RpcClient) SendRpcRequestWithTimeout(timeout time.Duration, rpcInvocation *RpcInvocation, responseMessage proto.Message) (*RpcDataPackage, error) {
	if c.Session == nil {
		return nil, ERR_NEED_INIT
	}

	now := time.Now().UnixNano()

	rpcDataPackage, err := rpcInvocation.GetRequestRpcDataPackage()
	if err != nil {
		return nil, err
	}

	ch := make(chan *RpcDataPackage)
	go c.asyncRequest(timeout, rpcDataPackage, ch)
	defer close(ch)
	// wait for message
	rsp := <-ch

	r := rsp
	if r == nil {
		return nil, ERR_RESPONSE_NIL //to ingore this nil value
	}

	errorCode := r.GetMeta().GetResponse().GetErrorCode()
	if errorCode > 0 {
		errMsg := fmt.Sprintf(LOG_SERVER_REPONSE_ERROR,
			errorCode, r.GetMeta().GetResponse().GetErrorText())
		return r, errors.New(errMsg)
	}

	response := r.GetData()
	if response != nil {
		err = proto.Unmarshal(response, responseMessage)
		if err != nil {
			return r, err
		}
	}

	took := TimetookInSeconds(now)
	Infof(LOG_CLIENT_TIMECOUST_INFO, *rpcInvocation.ServiceName, *rpcInvocation.MethodName, took)

	return r, nil

}
