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
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/funny/link"
	"github.com/jhunters/timewheel"
)

const REQUIRED_TYPE = "baidurpc.RpcDataPackage"

var (
	ERR_INVALID_TYPE        = errors.New("[codec-001]type mismatch, target type should be 'baidurpc.RpcDataPackage'")
	LOG_CLOSE_CONNECT_INFO  = "[codec-100]Do close connection. connection info:%v"
	chunkPackageCacheExpire = 60 * time.Second
)

/*
 Codec implements for RpcDataPackage.
*/
type RpcDataPackageCodec struct {
	readWriter io.ReadWriter
	closer     io.Closer
	p          *RpcDataPackageProtocol
	timeout    *time.Duration

	chunkPackageCache map[int64]*RpcDataPackage
}

// Here begin to implements link module Codec interface for RpcDataPackageCodec
/*
type Codec interface {
	Receive() (interface{}, error)
	Send(interface{}) error
	Close() error
	SetTimeout(timeout int)
}
*/

// send serialized data to target server by connection IO
// msg: param 'msg' must type of RpcDataPackage
func (r *RpcDataPackageCodec) Send(msg interface{}) error {

	if msg == nil {
		return errors.New("parameter 'msg' is nil")
	}

	v := reflect.ValueOf(msg)
	isPtr := false
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		isPtr = true
	}

	name := v.Type().String()
	if !strings.Contains(name, REQUIRED_TYPE) {
		return ERR_INVALID_TYPE
	}

	dataPackage := convertRpcDataPackage(msg, isPtr)

	rw := r.readWriter
	if r.timeout != nil {
		conn := rw.(net.Conn)
		conn.SetWriteDeadline(time.Now().Add(*r.timeout))
	}

	// check if use chunk mode
	chunkSize := dataPackage.chunkSize
	if chunkSize > 0 {
		dataPackageList := dataPackage.Chunk(int(chunkSize))
		for _, pack := range dataPackageList {
			err := pack.WriteIO(rw)
			if err != nil {
				return err
			}
		}
	} else {
		err := dataPackage.WriteIO(rw)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertRpcDataPackage(msg interface{}, isPtr bool) *RpcDataPackage {
	if isPtr {
		return msg.(*RpcDataPackage)
	}

	ret := msg.(RpcDataPackage)
	return &ret
}

// receive serialized data to target server by connection IO
// return param:
// 1. RpcDataPackage unserialized from connection io. or nil if exception found
// 2. a non-nil error if any io exception occured
func (r *RpcDataPackageCodec) Receive() (interface{}, error) {

	rw := r.readWriter

	if r.timeout != nil { // set time out
		conn := rw.(net.Conn)
		conn.SetReadDeadline(time.Now().Add(*r.timeout))
	}

	return r.doReceive(rw)

}

func (r *RpcDataPackageCodec) doReceive(conn io.ReadWriter) (interface{}, error) {
	dataPackage := NewRpcDataPackage()
	err := dataPackage.ReadIO(conn)
	if err != nil {
		if err == ERR_IGNORE_ERR {
			return nil, nil
		}
		return nil, err
	}

	// if chunk mode enabled
	if r.p.chunkSize > 0 {
		dataPackage.chunkSize = r.p.chunkSize
	}

	// check if chunk package
	if dataPackage.IsChunkPackage() {
		streamId := dataPackage.GetChunkStreamId()

		cachedPackage, exist := r.chunkPackageCache[streamId]
		if !exist {
			r.chunkPackageCache[streamId] = dataPackage
			cachedPackage = dataPackage

			// add task
			task := timewheel.Task{
				Data: streamId,
				TimeoutCallback: func(tt timewheel.Task) { // call back function on time out
					k := tt.Data.(int64)
					delete(r.chunkPackageCache, k)
				}}
			// add task and return unique task id
			r.p.tw.AddTask(chunkPackageCacheExpire, task) // add delay task

		} else {
			// if exist should merge data
			size := len(cachedPackage.Data) + len(dataPackage.Data)
			newData := make([]byte, size)
			copy(newData, cachedPackage.Data)
			copy(newData[len(cachedPackage.Data):], dataPackage.Data)
			cachedPackage.Data = newData
			r.chunkPackageCache[streamId] = cachedPackage
		}

		if dataPackage.IsFinalPackage() {
			delete(r.chunkPackageCache, streamId)
			// clear chunk status
			cachedPackage.ClearChunkStatus()
			return cachedPackage, nil
		} else {
			return r.doReceive(conn) // to receive next chunk package
		}
	}

	return dataPackage, nil
}

// do close connection io
// return non-nil if any error ocurred while doing close
func (r *RpcDataPackageCodec) Close() error {
	if r.closer != nil {
		log.Printf(LOG_CLOSE_CONNECT_INFO, r.closer)
		return r.closer.Close()
	}
	return nil
}

// set connection io read and write dead line
func (r *RpcDataPackageCodec) SetTimeout(timeout *time.Duration) {
	r.timeout = timeout
}

// Here begin to implements link module Protocol interface  for RpcDataPackageCodec
/*
type Protocol interface {
	NewCodec(rw io.ReadWriter) (Codec, error)
}

*/

// Protocol codec factory object for RpcDataPackage
type RpcDataPackageProtocol struct {
	timeout *time.Duration

	tw *timewheel.TimeWheel

	chunkSize uint32
}

// NewRpcDataPackageProtocol create a RpcDataPackageProtocol and start timewheel
func NewRpcDataPackageProtocol() (*RpcDataPackageProtocol, error) {
	protocol := &RpcDataPackageProtocol{}
	tw, err := timewheel.New(chunkExpireTimewheelInterval, uint16(chunkExpireTimeWheelSlot))
	if err != nil {
		return nil, err
	}
	protocol.tw = tw
	protocol.tw.Start()
	return protocol, nil
}

func (r *RpcDataPackageProtocol) NewCodec(rw io.ReadWriter) (link.Codec, error) {
	rpcDataPackage := &RpcDataPackageCodec{
		readWriter:        rw,
		p:                 r,
		timeout:           r.timeout,
		chunkPackageCache: make(map[int64]*RpcDataPackage),
	}

	rpcDataPackage.closer, _ = rw.(io.Closer)

	return rpcDataPackage, nil

}

// Stop
func (r *RpcDataPackageProtocol) Stop() {
	if r.tw != nil {
		r.tw.Stop()
	}
}

// Here end to implements link module Codec interface
