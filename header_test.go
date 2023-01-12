/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-04-26 18:18:59
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

	. "github.com/smartystreets/goconvey/convey"
)

// TestRpcDataWriteReader test head read and write
func TestRpcDataWriteReader(t *testing.T) {

	Convey("TestRpcDataWriteReader", t, func() {

		h := baidurpc.Header{}
		h.SetMagicCode([]byte("PRPB"))
		h.SetMessageSize(12300)
		h.SetMetaSize(59487)

		bs, _ := h.Write()
		So(len(bs), ShouldEqual, baidurpc.SIZE)

		h2 := baidurpc.Header{}
		h2.Read(bs)
		So(string(h.GetMagicCode()), ShouldEqual, string(h2.GetMagicCode()))
		So(h.GetMessageSize(), ShouldEqual, 12300)
		So(h.GetMetaSize(), ShouldEqual, 59487)
	})

}
