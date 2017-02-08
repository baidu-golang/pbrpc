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
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
)

// error log info definition
var ERR_NO_SNAPPY = errors.New("[marshal-002]Snappy compress not support yet.")
var ERR_IGNORE_ERR = errors.New("[marshal-001]Ingore error")
var ERR_META = errors.New("[marshal-003]Get nil value from Meta struct after marshal")

var LOG_INVALID_BYTES = "[marshal-004]Invalid byte array. maybe a broken byte stream. Received '%b'"

/*
 Data package for baidu RPC.
 all request and response data package should apply this.

-----------------------------------
| Head | Meta | Data | Attachment |
-----------------------------------

1. <Head> with fixed 12 byte length as follow format
----------------------------------------------
| PRPC | MessageSize(int32) | MetaSize(int32) |
----------------------------------------------
MessageSize = totalSize - 12(Fixed Head Size)
MetaSize = Meta object size

2. <Meta> body proto description as follow
message RpcMeta {
    optional RpcRequestMeta request = 1;
    optional RpcResponseMeta response = 2;
    optional int32 compress_type = 3; // 0:nocompress 1:Snappy 2:gzip
    optional int64 correlation_id = 4;
    optional int32 attachment_size = 5;
    optional ChunkInfo chuck_info = 6;
    optional bytes authentication_data = 7;
};

message Request {
    required string service_name = 1;
    required string method_name = 2;
    optional int64 log_id = 3;
};

message Response {
    optional int32 error_code = 1;
    optional string error_text = 2;
};

messsage ChunkInfo {
        required int64 stream_id = 1;
        required int64 chunk_id = 2;
};

3. <Data> customize transport data message.

4. <Attachment> attachment body data message

*/
type RpcDataPackage struct {
	Head       *Header
	Meta       *RpcMeta
	Data       []byte
	Attachment []byte
}

func NewRpcDataPackage() *RpcDataPackage {
	data := RpcDataPackage{}
	doInit(&data)

	data.GetMeta().Response = &Response{}

	return &data
}

func (r *RpcDataPackage) MagicCode(magicCode string) {
	if len(magicCode) != 4 {
		return
	}

	initHeader(r)
	r.Head.SetMagicCode([]byte(magicCode))

}

func (r *RpcDataPackage) GetMagicCode() string {
	initHeader(r)
	return string(r.Head.GetMagicCode())

}

func initHeader(r *RpcDataPackage) {
	if r.Head == nil {
		r.Head = &Header{}
	}
}

func initRpcMeta(r *RpcDataPackage) {
	if r.Meta == nil {
		r.Meta = &RpcMeta{}
	}
}

func initRequest(r *RpcDataPackage) {
	initRpcMeta(r)

	request := r.Meta.Request
	if request == nil {
		r.Meta.Request = &Request{}
	}

}

func initResponse(r *RpcDataPackage) {
	initRpcMeta(r)

	response := r.Meta.Response
	if response == nil {
		r.Meta.Response = &Response{}
	}

}

func (r *RpcDataPackage) ServiceName(serviceName string) *RpcDataPackage {
	initRequest(r)

	r.Meta.Request.ServiceName = &serviceName

	return r
}

func (r *RpcDataPackage) MethodName(methodName string) *RpcDataPackage {
	initRequest(r)

	r.Meta.Request.MethodName = &methodName

	return r
}

func (r *RpcDataPackage) SetData(Data []byte) *RpcDataPackage {
	r.Data = Data
	return r
}

func (r *RpcDataPackage) SetAttachment(Attachment []byte) *RpcDataPackage {
	r.Attachment = Attachment
	return r
}

func (r *RpcDataPackage) AuthenticationData(authenticationData []byte) *RpcDataPackage {
	initRpcMeta(r)

	r.Meta.AuthenticationData = authenticationData
	return r
}

func (r *RpcDataPackage) CorrelationId(correlationId int64) *RpcDataPackage {
	initRpcMeta(r)

	r.Meta.CorrelationId = &correlationId
	return r
}

func (r *RpcDataPackage) CompressType(compressType int32) *RpcDataPackage {
	initRpcMeta(r)

	r.Meta.CompressType = &compressType
	return r
}

func (r *RpcDataPackage) LogId(logId int64) *RpcDataPackage {
	initRequest(r)

	r.Meta.Request.LogId = &logId

	return r
}

func (r *RpcDataPackage) GetLogId() int64 {
	initRequest(r)
	return r.Meta.Request.GetLogId()
}

func (r *RpcDataPackage) ErrorCode(errorCode int32) *RpcDataPackage {
	initResponse(r)

	r.Meta.Response.ErrorCode = &errorCode

	return r
}

func (r *RpcDataPackage) ErrorText(errorText string) *RpcDataPackage {
	initResponse(r)

	r.Meta.Response.ErrorText = &errorText

	return r
}

func (r *RpcDataPackage) ExtraParams(extraParams []byte) *RpcDataPackage {
	initRequest(r)

	r.Meta.Request.ExtraParam = extraParams

	return r
}

func (r *RpcDataPackage) ChunkInfo(streamId int64, chunkId int64) *RpcDataPackage {
	chunkInfo := ChunkInfo{}
	chunkInfo.StreamId = &streamId
	chunkInfo.ChunkId = &chunkId
	initRpcMeta(r)
	r.Meta.ChunkInfo = &chunkInfo
	return r
}

func doInit(r *RpcDataPackage) {
	initHeader(r)
	initRequest(r)
	initResponse(r)
}

func (r *RpcDataPackage) GetHead() *Header {
	if r.Head == nil {
		return nil
	}
	return r.Head
}

func (r *RpcDataPackage) GetMeta() *RpcMeta {
	if r.Meta == nil {
		return nil
	}
	return r.Meta
}

func (r *RpcDataPackage) GetData() []byte {
	return r.Data
}

func (r *RpcDataPackage) GetAttachment() []byte {
	return r.Attachment
}

/*
  Convert RpcPackage to byte array
*/
func (r *RpcDataPackage) WriteIO(rw io.ReadWriter) error {

	bytes, err := r.Write()
	if err != nil {
		return err
	}

	_, err = rw.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

/*
  Convert RpcPackage to byte array
*/
func (r *RpcDataPackage) Write() ([]byte, error) {
	doInit(r)

	var totalSize int32 = 0
	var dataSize int32 = 0
	var err error
	if r.Data != nil {
		compressType := r.GetMeta().GetCompressType()
		if compressType == COMPRESS_GZIP {
			r.Data, err = GZIP(r.Data)
			if err != nil {
				return nil, err
			}
		} else if compressType == COMPRESS_SNAPPY {
			dst := make([]byte, snappy.MaxEncodedLen(len(r.Data)))
			r.Data = snappy.Encode(dst, r.Data)
		}

		dataSize = int32(len(r.Data))
		totalSize = totalSize + dataSize
	}

	var attachmentSize int32 = 0
	if r.Attachment != nil {
		attachmentSize = int32(len(r.Attachment))
		totalSize = totalSize + attachmentSize
	}

	r.Meta.AttachmentSize = proto.Int32(int32(attachmentSize))

	metaBytes, err := proto.Marshal(r.Meta)
	if err != nil {
		return nil, err
	}

	if metaBytes == nil {
		return nil, ERR_META
	}

	rpcMetaSize := int32(len(metaBytes))
	totalSize = totalSize + rpcMetaSize

	r.Head.SetMetaSize(int32(rpcMetaSize))
	r.Head.SetMessageSize(int32(totalSize)) // set message body size

	r.Head.MessageSize = int32(totalSize)
	r.Head.MetaSize = int32(rpcMetaSize)
	buf := new(bytes.Buffer)

	headBytes, _ := r.Head.Write()
	binary.Write(buf, binary.BigEndian, headBytes)
	binary.Write(buf, binary.BigEndian, metaBytes)

	if r.Data != nil {
		binary.Write(buf, binary.BigEndian, r.Data)
	}

	if r.Attachment != nil {
		binary.Write(buf, binary.BigEndian, r.Attachment)
	}

	return buf.Bytes(), nil
}

func checkSize(current, expect int, bytes []byte) error {
	if current < expect {
		message := fmt.Sprintf(LOG_INVALID_BYTES, bytes)
		return errors.New(message)

	}
	return nil
}

/*
Read byte array and initialize RpcPackage
*/
func (r *RpcDataPackage) ReadIO(rw io.ReadWriter) error {
	if rw == nil {
		return errors.New("bytes is nil")
	}

	doInit(r)

	// read Head
	head := make([]byte, SIZE)
	_, err := rw.Read(head)
	if err != nil {
		log.Println("Read head error", err)
		// only to close current connection
		return ERR_IGNORE_ERR
	}

	// unmarshal Head message
	r.Head.Read(head)

	// get RPC Meta size
	metaSize := r.Head.GetMetaSize()
	totalSize := r.Head.GetMessageSize()
	if totalSize <= 0 {
		// maybe heart beat data message, so do ignore here
		return ERR_IGNORE_ERR
	}

	// read left
	leftSize := totalSize
	body := make([]byte, leftSize)

	rw.Read(body)

	proto.Unmarshal(body[0:metaSize], r.Meta)

	attachmentSize := r.Meta.GetAttachmentSize()
	dataSize := leftSize - metaSize - attachmentSize

	dataOffset := metaSize
	if dataSize > 0 {
		dataOffset = dataSize + metaSize
		r.Data = body[metaSize:dataOffset]

		compressType := r.GetMeta().GetCompressType()
		if compressType == COMPRESS_GZIP {
			r.Data, err = GUNZIP(r.Data)

			if err != nil {
				return err
			}
		} else if compressType == COMPRESS_SNAPPY {
			dst := make([]byte, 1)
			r.Data, err = snappy.Decode(dst, r.Data)
			if err != nil {
				return err
			}
		}
	}
	// if need read Attachment
	if attachmentSize > 0 {
		r.Attachment = body[dataOffset:leftSize]
	}

	return nil
}

func (r *RpcDataPackage) Read(b []byte) error {

	if b == nil {
		return errors.New("b is nil")
	}

	buf := bytes.NewBuffer(b)

	return r.ReadIO(buf)

}
