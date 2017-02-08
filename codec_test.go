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
package pbrpc_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	pbrpc "github.com/baidu-golang/pbrpc"
)

func TestCodecSend(t *testing.T) {

	protocol := &pbrpc.RpcDataPackageProtocol{}

	buf := new(bytes.Buffer)

	rpcDataPackagecodec, _ := protocol.NewCodec(buf)

	rpcDataPackage := initRpcDataPackage()

	rpcDataPackagecodec.Send(rpcDataPackage)

	r2 := pbrpc.RpcDataPackage{}
	err := r2.ReadIO(buf)

	if err != nil {
		t.Error(err)
	}

	err = equalRpcDataPackage(r2)
	if err != nil {
		t.Error(err)
	}
}

func TestCodecReceive(t *testing.T) {

	protocol := &pbrpc.RpcDataPackageProtocol{}

	buf := new(bytes.Buffer)
	// write prepare value to buf
	rpcDataPackage := initRpcDataPackage()
	bytes, err := rpcDataPackage.Write()
	if err != nil {
		t.Error(err)
	}

	buf.Write(bytes)
	rpcDataPackagecodec, _ := protocol.NewCodec(buf)

	rsp, err := rpcDataPackagecodec.Receive()

	mc := rsp.(*pbrpc.RpcDataPackage).GetMagicCode()
	if !strings.EqualFold(magicCode, mc) {
		t.Error(fmt.Sprintf("expect magic code is '%s' but actual is '%s'", magicCode, mc))
	}

}
