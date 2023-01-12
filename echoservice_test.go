/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-03 14:08:57
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
/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-03 14:08:57
 */
package baidurpc_test

import (
	"context"
	"fmt"
	"time"

	baidurpc "github.com/baidu-golang/pbrpc"
)

// EchoService define service
type EchoService struct {
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) Echo(c context.Context, in *DataMessage) (*DataMessage, context.Context) {
	var ret = "hello "

	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := DataMessage{Name: ret}

	// time.Sleep(50 * time.Second)

	// bind with err
	// cc = baidurpc.BindError(cc, errors.New("manule error"))
	return &dm, context.TODO()
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) EchoWithAttchement(c context.Context, in *DataMessage) (*DataMessage, context.Context) {
	var ret = "hello "

	// get attchement
	attachement := baidurpc.Attachement(c)

	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := DataMessage{Name: ret}

	att := "I am a attachement, " + string(attachement)

	// bind attachment
	cc := baidurpc.BindAttachement(context.Background(), []byte(att))
	return &dm, cc
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) EchoWithCustomizedError(c context.Context, in *DataMessage) (*DataMessage, context.Context) {
	var ret = "hello "

	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := DataMessage{Name: ret}

	// bind with err
	cc := baidurpc.BindError(context.Background(), fmt.Errorf("this is customized error"))
	return &dm, cc
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) EchoWithoutContext(c context.Context, in *DataMessage) *DataMessage {
	var ret = "hello "

	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := DataMessage{Name: ret}

	return &dm
}

// Echo  test publish method with return type has context argument
func (rpc *EchoService) EchoSlowTest(c context.Context, in *DataMessage) *DataMessage {
	var ret = "hello "

	time.Sleep(2 * time.Second)
	if len(in.Name) == 0 {
		ret = ret + "veryone"
	} else {
		ret = ret + in.Name
	}

	// return result
	dm := DataMessage{Name: ret}

	return &dm
}
