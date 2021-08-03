/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-03 14:08:57
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
/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-03 14:08:57
 */
package baidurpc

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
)

//手工定义pb生成的代码, tag 格式 = protobuf:"type,order,req|opt|rep|packed,name=fieldname"
type EchoMessage struct {
	Name string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
}

func (m *EchoMessage) Reset()         { *m = EchoMessage{} }
func (m *EchoMessage) String() string { return proto.CompactTextString(m) }
func (*EchoMessage) ProtoMessage()    {}

func (m *EchoMessage) GetName() string {
	return m.Name
}

// EchoService define service
type EchoService struct {
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) Echo(c context.Context, in *EchoMessage) (*EchoMessage, context.Context) {
	var ret = "hello "

	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := EchoMessage{ret}

	// bind with err
	// cc = baidurpc.BindError(cc, errors.New("manule error"))
	return &dm, context.TODO()
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) EchoWithAttchement(c context.Context, in *EchoMessage) (*EchoMessage, context.Context) {
	var ret = "hello "

	// get attchement
	attachement := Attachement(c)

	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := EchoMessage{ret}

	att := "I am a attachement" + string(attachement)

	// bind attachment
	cc := BindAttachement(context.Background(), []byte(att))
	return &dm, cc
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) EchoWithCustomizedError(c context.Context, in *EchoMessage) (*EchoMessage, context.Context) {
	var ret = "hello "

	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := EchoMessage{ret}

	// bind with err
	cc := BindError(context.Background(), fmt.Errorf("this is customized error"))
	return &dm, cc
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) EchoWithoutContext(c context.Context, in *EchoMessage) *EchoMessage {
	var ret = "hello "

	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := EchoMessage{ret}

	return &dm
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) EchoSlowTest(c context.Context, in *EchoMessage) *EchoMessage {
	var ret = "hello "

	time.Sleep(2 * time.Second)
	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := EchoMessage{ret}

	return &dm
}
