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
package baidurpc

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"github.com/jhunters/timewheel"
)

// HaRpcClient high avialbe RpcClient
type HaRpcClient struct {
	rpcClients []*RpcClient

	current int32

	locker sync.Mutex

	tw *timewheel.TimeWheel
}

// NewBatchTCPConnection to create batch connection
func NewBatchTCPConnection(urls []URL, timeout time.Duration) ([]Connection, error) {
	if len(urls) == 0 {
		return nil, errors.New("param 'urls' is empty")
	}

	var result []Connection
	for _, url := range urls {
		conn, err := NewTCPConnection(url, &timeout)
		if err != nil {
			log.Println("create connection failed:", err)
			continue
		}
		result = append(result, conn)
	}

	return result, nil
}

// CloseBatchConnection close batch connections
func CloseBatchConnection(connections []Connection) {
	for _, conn := range connections {
		conn.Close()
	}
}

// NewHaRpcCient
func NewHaRpcCient(connections []Connection) (*HaRpcClient, error) {
	return NewHaRpcCientWithTimewheelSetting(connections, defaultTimewheelInterval, uint16(defaultTimewheelSlot))
}

// NewHaRpcCient
func NewHaRpcCientWithTimewheelSetting(connections []Connection, timewheelInterval time.Duration, timewheelSlot uint16) (*HaRpcClient, error) {
	if len(connections) == 0 {
		return nil, errors.New("param 'connections' is empty")
	}

	rpcClients := make([]*RpcClient, len(connections))
	for idx, connection := range connections {
		rpcClient, err := NewRpcCient(connection)
		if err != nil {
			return nil, err
		}
		rpcClients[idx] = rpcClient
	}
	result := &HaRpcClient{rpcClients: rpcClients, current: 0}
	result.tw, _ = timewheel.New(timewheelInterval, timewheelSlot)
	result.tw.Start()
	return result, nil
}

// electClient
func (c *HaRpcClient) electClient() *RpcClient {
	if len(c.rpcClients) == 0 {
		return nil
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	client := c.rpcClients[int(c.current)%len(c.rpcClients)]
	c.current++
	return client
}

// SendRpcRequest send rpc request by elect one client
func (c *HaRpcClient) SendRpcRequest(rpcInvocation *RpcInvocation, responseMessage proto.Message) (*RpcDataPackage, error) {
	var errRet error
	size := len(c.rpcClients)
	if size == 0 {
		return nil, errors.New("no rpc client avaible")
	}

	for i := 0; i < size; i++ {
		rpcClient := c.electClient()
		if rpcClient == nil {
			return nil, errors.New("no rpc client avaible")
		}

		data, err := rpcClient.SendRpcRequest(rpcInvocation, responseMessage)
		if err == nil {
			return data, nil
		} else {
			errRet = err
		}
	}

	return nil, errRet
}

// SendRpcRequest send rpc request by elect one client with timeout feature
func (c *HaRpcClient) SendRpcRequestWithTimeout(timeout time.Duration, rpcInvocation *RpcInvocation, responseMessage proto.Message) (*RpcDataPackage, error) {
	var errRet error
	size := len(c.rpcClients)
	if size == 0 {
		return nil, errors.New("no rpc client avaible")
	}

	ch := make(chan *RpcDataPackage)
	go c.asyncRequest(timeout, rpcInvocation, responseMessage, ch)
	defer close(ch)
	// wait for message
	rsp := <-ch

	return rsp, errRet
}

// asyncRequest
func (c *HaRpcClient) asyncRequest(timeout time.Duration, rpcInvocation *RpcInvocation, responseMessage proto.Message, ch chan<- *RpcDataPackage) {
	request, err := rpcInvocation.GetRequestRpcDataPackage()
	if err != nil {
		errorcode := int32(ST_ERROR)
		request.ErrorCode(errorcode)
		errormsg := err.Error()
		request.ErrorText(errormsg)

		ch <- request
		return
	}

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

	rsp, err := c.SendRpcRequest(rpcInvocation, responseMessage)
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

// Close do close all client
func (c *HaRpcClient) Close() {
	if c.tw != nil {
		c.tw.Stop()
	}
	for _, client := range c.rpcClients {
		client.Close()
	}
}
