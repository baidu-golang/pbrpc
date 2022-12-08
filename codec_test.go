/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-07-24 16:54:14
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
package baidurpc_test

import (
	"bytes"
	"testing"

	baidurpc "github.com/baidu-golang/pbrpc"
	. "github.com/smartystreets/goconvey/convey"
)

// TestCodecSend
func TestCodecSend(t *testing.T) {
	Convey("TestCodecSend", t, func() {
		protocol := &baidurpc.RpcDataPackageProtocol[*baidurpc.RpcDataPackage, *baidurpc.RpcDataPackage]{}

		buf := new(bytes.Buffer)

		rpcDataPackagecodec, _ := protocol.NewCodec(buf)

		rpcDataPackage := initRpcDataPackage()

		rpcDataPackagecodec.Send(rpcDataPackage)

		r2 := baidurpc.RpcDataPackage{}
		err := r2.ReadIO(buf)
		So(err, ShouldBeNil)

		err = equalRpcDataPackage(r2, t)
		So(err, ShouldBeNil)
	})

}

// TestCodecReceive
func TestCodecReceive(t *testing.T) {

	Convey("TestCodecReceive", t, func() {
		protocol := &baidurpc.RpcDataPackageProtocol[*baidurpc.RpcDataPackage, *baidurpc.RpcDataPackage]{}

		buf := new(bytes.Buffer)
		// write prepare value to buf
		rpcDataPackage := initRpcDataPackage()
		bytes, err := rpcDataPackage.Write()
		So(err, ShouldBeNil)

		buf.Write(bytes)
		rpcDataPackagecodec, _ := protocol.NewCodec(buf)

		rsp, err := rpcDataPackagecodec.Receive()
		So(err, ShouldBeNil)
		mc := rsp.GetMagicCode()
		So(string(magicCode), ShouldEqual, string(mc))
	})

}
