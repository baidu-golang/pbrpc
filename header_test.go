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
package baidurpc_test

import (
	"bytes"
	"testing"

	baidurpc "github.com/baidu-golang/pbrpc"
)

func TestRpcDataWriteReader(t *testing.T) {

	h := baidurpc.Header{}
	h.SetMagicCode([]byte("PRPB"))
	h.SetMessageSize(12300)
	h.SetMetaSize(59487)

	bs, _ := h.Write()

	if len(bs) != baidurpc.SIZE {
		t.Errorf("current head size is '%d', should be '%d'", len(bs), baidurpc.SIZE)
	}

	h2 := baidurpc.Header{}
	h2.Read(bs)
	if !bytes.Equal(h.GetMagicCode(), h2.GetMagicCode()) {
		t.Errorf("magic code is not same. expect '%b' actual is '%b'", h.GetMagicCode(), h2.GetMagicCode())
	}

}
