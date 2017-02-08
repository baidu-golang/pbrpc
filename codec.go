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
	"io"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/baidu-golang/pbrpc/link"
	"github.com/golang/glog"
)

const REQUIRED_TYPE = "pbrpc.RpcDataPackage"

var ERR_INVALID_TYPE = errors.New("[codec-001]type mismatch, target type should be 'pbrpc.RpcDataPackage'")
var LOG_CLOSE_CONNECT_INFO = "[codec-100]Do close connection. connection info:%v"

/*
 Codec implements for RpcDataPackage.
*/
type RpcDataPackageCodec struct {
	readWriter io.ReadWriter
	closer     io.Closer
	p          *RpcDataPackageProtocol
	timeout    *int
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
		t := *r.timeout
		timeout := time.Duration(t) * time.Second
		conn.SetWriteDeadline(time.Now().Add(timeout))
	}

	err := dataPackage.WriteIO(rw)
	if err != nil {
		return err
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
		t := *r.timeout
		timeout := time.Duration(t) * time.Second
		conn.SetReadDeadline(time.Now().Add(timeout))
	}

	dataPackage := NewRpcDataPackage()
	err := dataPackage.ReadIO(rw)
	if err != nil {
		if err == ERR_IGNORE_ERR {
			return nil, nil
		}
		return nil, err
	}

	return dataPackage, nil
}

// do close connection io
// return non-nil if any error ocurred while doing close
func (r *RpcDataPackageCodec) Close() error {
	if r.closer != nil {
		glog.Infof(LOG_CLOSE_CONNECT_INFO, r.closer)
		return r.closer.Close()
	}
	return nil
}

// set connection io read and write dead line
func (r *RpcDataPackageCodec) SetTimeout(timeout *int) {
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
}

func (r *RpcDataPackageProtocol) NewCodec(rw io.ReadWriter) (link.Codec, error) {
	rpcDataPackage := &RpcDataPackageCodec{
		readWriter: rw,
		p:          r,
	}

	rpcDataPackage.closer, _ = rw.(io.Closer)

	return rpcDataPackage, nil

}

// Here end to implements link module Codec interface
