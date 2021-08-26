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
	"testing"

	baidurpc "github.com/baidu-golang/pbrpc"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
)

var sericeName = "thisIsAServiceName"
var methodName = "thisIsAMethodName"
var magicCode = "PRPC"
var logId int64 = 1001
var correlationId int64 = 20001
var data []byte = []byte{1, 2, 3, 1, 2, 3, 1, 1, 2, 2, 20}
var attachment []byte = []byte{2, 2, 2, 2, 2, 1, 1, 1, 1}

func initRpcDataPackage() *baidurpc.RpcDataPackage {

	rpcDataPackage := baidurpc.RpcDataPackage{}

	rpcDataPackage.MagicCode(magicCode)
	rpcDataPackage.SetData(data)
	rpcDataPackage.ServiceName(sericeName)
	rpcDataPackage.MethodName(methodName)

	rpcDataPackage.LogId(logId)
	rpcDataPackage.CorrelationId(correlationId)

	rpcDataPackage.SetAttachment(attachment)

	return &rpcDataPackage
}

func equalRpcDataPackage(r baidurpc.RpcDataPackage, t *testing.T) error {
	Convey("test RpcDataPackage", func() {
		So(sericeName, ShouldEqual, r.Meta.Request.ServiceName)
		So(methodName, ShouldEqual, r.Meta.Request.MethodName)
		So(string(magicCode), ShouldEqual, string(r.GetMagicCode()))
		So(r.Meta.Request.LogId, ShouldEqual, logId)
		So(r.Meta.CorrelationId, ShouldEqual, correlationId)
		So(len(data), ShouldEqual, len(r.Data))
		So(len(attachment), ShouldEqual, len(r.Attachment))
	})
	return nil
}

func validateRpcDataPackage(t *testing.T, r2 baidurpc.RpcDataPackage) {
	Convey("validateRpcDataPackage", func() {
		So(string(magicCode), ShouldEqual, string(r2.GetMagicCode()))
		So(sericeName, ShouldEqual, r2.GetMeta().GetRequest().GetServiceName())
		So(methodName, ShouldEqual, r2.GetMeta().GetRequest().GetMethodName())
	})

}

// TestWriteReaderWithMockData
func TestWriteReaderWithMockData(t *testing.T) {

	Convey("TestWriteReaderWithMockData", t, func() {
		rpcDataPackage := initRpcDataPackage()

		b, err := rpcDataPackage.Write()
		if err != nil {
			t.Error(err.Error())
		}

		r2 := baidurpc.RpcDataPackage{}

		err = r2.Read(b)
		if err != nil {
			t.Error(err.Error())
		}

		validateRpcDataPackage(t, r2)
	})

}

// WriteReaderWithRealData
func WriteReaderWithRealData(rpcDataPackage *baidurpc.RpcDataPackage,
	compressType int32, t *testing.T) {

	Convey("Test with real data", func() {
		dataMessage := DataMessage{}
		name := "hello, xiemalin. this is repeated string aaaaaaaaaaaaaaaaaaaaaa"
		dataMessage.Name = name

		data, err := proto.Marshal(&dataMessage)
		So(err, ShouldBeNil)
		rpcDataPackage.SetData(data)

		b, err := rpcDataPackage.Write()
		So(err, ShouldBeNil)

		r2 := baidurpc.RpcDataPackage{}
		r2.CompressType(compressType)

		err = r2.Read(b)
		So(err, ShouldBeNil)

		validateRpcDataPackage(t, r2)

		newData := r2.GetData()
		dataMessage2 := DataMessage{}
		proto.Unmarshal(newData, &dataMessage2)

		So(name, ShouldEqual, dataMessage2.Name)
	})

}

// TestWriteReaderWithRealData
func TestWriteReaderWithRealData(t *testing.T) {

	Convey("TestWriteReaderWithRealData", t, func() {
		rpcDataPackage := initRpcDataPackage()
		WriteReaderWithRealData(rpcDataPackage, baidurpc.COMPRESS_NO, t)
	})

}

// TestWriteReaderWithGZIP
func TestWriteReaderWithGZIP(t *testing.T) {

	Convey("TestWriteReaderWithGZIP", t, func() {
		rpcDataPackage := initRpcDataPackage()
		rpcDataPackage.CompressType(baidurpc.COMPRESS_GZIP)
		WriteReaderWithRealData(rpcDataPackage, baidurpc.COMPRESS_GZIP, t)
	})

}

// TestWriteReaderWithSNAPPY
func TestWriteReaderWithSNAPPY(t *testing.T) {

	Convey("TestWriteReaderWithSNAPPY", t, func() {
		rpcDataPackage := initRpcDataPackage()

		rpcDataPackage.CompressType(baidurpc.COMPRESS_SNAPPY)

		WriteReaderWithRealData(rpcDataPackage, baidurpc.COMPRESS_SNAPPY, t)
	})

}

// TestChunk
func TestChunk(t *testing.T) {

	Convey("Test package chunk", t, func() {
		rpcDataPackage := initRpcDataPackage()

		Convey("Test package chunk with invalid chunk size", func() {
			chunkSize := 0
			chunkPackages := rpcDataPackage.Chunk(chunkSize)
			So(len(chunkPackages), ShouldEqual, 1)
		})

		Convey("Test package chunk with chunk size", func() {
			chunkSize := 1
			count := len(rpcDataPackage.Data)
			chunkPackages := rpcDataPackage.Chunk(chunkSize)
			So(len(chunkPackages), ShouldEqual, count)
			So(len(chunkPackages[0].Data), ShouldEqual, 1)
		})

		Convey("Test package chunk with large chunk size", func() {
			chunkSize := 100
			datasize := len(rpcDataPackage.Data)
			chunkPackages := rpcDataPackage.Chunk(chunkSize)
			So(len(chunkPackages), ShouldEqual, 1)
			So(len(chunkPackages[0].Data), ShouldEqual, datasize)
		})
	})

}
