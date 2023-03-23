// Go support for Protocol Buffers RPC which compatible with https://github.com/Baidu-ecom/Jprotobuf-rpc-socket
//
// Copyright 2002-2007 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package baidurpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"

	"github.com/golang/snappy"
	"google.golang.org/protobuf/proto"
)

// error log info definition
var (
	errIgnoreErr      = errors.New("[marshal-001]Ingore error")
	errMeta           = errors.New("[marshal-003]Get nil value from Meta struct after marshal")
	LOG_INVALID_BYTES = "[marshal-004]Invalid byte array. maybe a broken byte stream. Received '%b'"
)

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
	    optional Request request = 1;
	    optional Response response = 2;
	    optional int32 compress_type = 3; // 0:nocompress 1:Snappy 2:gzip
	    optional int64 correlation_id = 4;
	    optional int32 attachment_size = 5;
	    optional ChuckInfo chuck_info = 6;
	    optional bytes authentication_data = 7;
	};

	message Request {
	    required string service_name = 1;
	    required string method_name = 2;
	    optional int64 log_id = 3;
		optional int64 traceId=4;
		optional int64 spanId=5;
		optional int64 parentSpanId=6;
		repeat RpcRequestMetaExtField extFields = 7;
	};

	message RpcRequestMetaExtField {
		optional string key = 1;
		optional string value = 2;
	}

	message Response {
	    optional int32 error_code = 1;
	    optional string error_text = 2;
	};

	messsage ChuckInfo {
	    required int64 stream_id = 1;
	    required int64 chunk_id = 2;
	};

3. <Data> customize transport data message.

4. <Attachment> attachment body data message
*/
type RpcDataPackage struct {
	Head       *Header  // rpc head
	Meta       *RpcMeta // rpc meta
	Data       []byte
	Attachment []byte

	// private field
	chunkSize uint32
}

// NewRpcDataPackage returns  a new RpcDataPackage and init all fields
func NewRpcDataPackage() *RpcDataPackage {
	data := RpcDataPackage{}
	doInit(&data)

	data.GetMeta().Response = &Response{}

	return &data
}

// Clear to clear and init all fields
func (r *RpcDataPackage) Clear() {
	// r.Head = &Header{}
	// r.Meta = &RpcMeta{}
	request := r.Meta.Request
	if request == nil {
		r.Meta.Request = &Request{}
	}
	response := r.Meta.Response
	if response == nil {
		r.Meta.Response = &Response{}
	}
	r.Data = nil
	r.Attachment = nil
	r.ClearChunkStatus()
}

// MagicCode set magic code field
func (r *RpcDataPackage) MagicCode(magicCode string) {
	if len(magicCode) != 4 {
		return
	}

	initHeader(r)
	r.Head.SetMagicCode([]byte(magicCode))

}

// GetMagicCode return magic code value
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

func initChuckInfo(r *RpcDataPackage) {
	initRpcMeta(r)
	if r.Meta.ChuckInfo == nil {
		r.Meta.ChuckInfo = &ChunkInfo{}
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

// ServiceName set service name field
func (r *RpcDataPackage) ServiceName(serviceName string) *RpcDataPackage {
	initRequest(r)

	r.Meta.Request.ServiceName = serviceName

	return r
}

// MethodName set method name field
func (r *RpcDataPackage) MethodName(methodName string) *RpcDataPackage {
	initRequest(r)

	r.Meta.Request.MethodName = methodName

	return r
}

// SetData set data
func (r *RpcDataPackage) SetData(Data []byte) *RpcDataPackage {
	r.Data = Data
	return r
}

// SetAttachment set attachment
func (r *RpcDataPackage) SetAttachment(Attachment []byte) *RpcDataPackage {
	r.Attachment = Attachment
	return r
}

// AuthenticationData set authentication data
func (r *RpcDataPackage) AuthenticationData(authenticationData []byte) *RpcDataPackage {
	initRpcMeta(r)

	r.Meta.AuthenticationData = authenticationData
	return r
}

// CorrelationId set correlationId data
func (r *RpcDataPackage) CorrelationId(correlationId int64) *RpcDataPackage {
	initRpcMeta(r)

	r.Meta.CorrelationId = correlationId
	return r
}

// CompressType set compress type data
func (r *RpcDataPackage) CompressType(compressType int32) *RpcDataPackage {
	initRpcMeta(r)

	r.Meta.CompressType = compressType
	return r
}

// LogId set log id data
func (r *RpcDataPackage) LogId(logId int64) *RpcDataPackage {
	initRequest(r)

	r.Meta.Request.LogId = logId

	return r
}

// GetLogId return log id
func (r *RpcDataPackage) GetLogId() int64 {
	initRequest(r)
	return r.Meta.Request.GetLogId()
}

// TraceId set trace id data
func (r *RpcDataPackage) TraceId(traceId int64) *RpcDataPackage {
	initRequest(r)
	r.Meta.Request.TraceId = traceId
	return r
}

// GetTraceId return trace id
func (r *RpcDataPackage) GetTraceId() int64 {
	initRequest(r)
	return r.Meta.Request.TraceId
}

// SpanId set span id
func (r *RpcDataPackage) SpanId(spanId int64) *RpcDataPackage {
	initRequest(r)
	r.Meta.Request.SpanId = spanId
	return r
}

// GetSpanId return span id
func (r *RpcDataPackage) GetSpanId() int64 {
	initRequest(r)
	return r.Meta.Request.SpanId
}

// ParentSpanId set parent span id
func (r *RpcDataPackage) ParentSpanId(parentSpanId int64) *RpcDataPackage {
	initRequest(r)
	r.Meta.Request.ParentSpanId = parentSpanId
	return r
}

// GetParentSpanId return parent span id
func (r *RpcDataPackage) GetParentSpanId() int64 {
	initRequest(r)
	return r.Meta.Request.ParentSpanId
}

// RpcRequestMetaExt set rpc request meta extendsion fields
func (r *RpcDataPackage) RpcRequestMetaExt(ext map[string]string) *RpcDataPackage {
	initRequest(r)
	extMap := make([]*RpcRequestMetaExtField, 0)
	for key, value := range ext {
		extfield := &RpcRequestMetaExtField{Key: key, Value: value}
		extMap = append(extMap, extfield)
	}
	r.Meta.Request.RpcRequestMetaExt = extMap
	return r
}

// GetRpcRequestMetaExt return rpc request meta extendstion
func (r *RpcDataPackage) GetRpcRequestMetaExt() map[string]string {
	initRequest(r)
	ret := make(map[string]string)
	for _, rr := range r.Meta.Request.RpcRequestMetaExt {
		ret[rr.Key] = rr.Value
	}

	return ret
}

// ErrorCode set error code field
func (r *RpcDataPackage) ErrorCode(errorCode int32) *RpcDataPackage {
	initResponse(r)

	r.Meta.Response.ErrorCode = errorCode

	return r
}

// ErrorText set error text field
func (r *RpcDataPackage) ErrorText(errorText string) *RpcDataPackage {
	initResponse(r)

	r.Meta.Response.ErrorText = errorText

	return r
}

// ExtraParams set extra parameters field
func (r *RpcDataPackage) ExtraParams(extraParams []byte) *RpcDataPackage {
	initRequest(r)

	r.Meta.Request.ExtraParam = extraParams

	return r
}

// ChuckInfo set chuck info
func (r *RpcDataPackage) ChuckInfo(streamId int64, chunkId int64) *RpcDataPackage {
	ChuckInfo := ChunkInfo{}
	ChuckInfo.StreamId = streamId
	ChuckInfo.ChunkId = chunkId
	initRpcMeta(r)
	r.Meta.ChuckInfo = &ChuckInfo
	return r
}

func doInit(r *RpcDataPackage) {
	initHeader(r)
	initRequest(r)
	initResponse(r)
}

// GetHead return Header data
func (r *RpcDataPackage) GetHead() *Header {
	if r.Head == nil {
		return nil
	}
	return r.Head
}

// GetMeta return RpcMeta data
func (r *RpcDataPackage) GetMeta() *RpcMeta {
	if r.Meta == nil {
		return nil
	}
	return r.Meta
}

// GetData return data field
func (r *RpcDataPackage) GetData() []byte {
	return r.Data
}

// GetAttachment return attachment field
func (r *RpcDataPackage) GetAttachment() []byte {
	return r.Attachment
}

/*
Convert RpcPackage to byte array
*/
func (r *RpcDataPackage) WriteIO(rw io.Writer) error {

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

	r.Meta.AttachmentSize = int32(attachmentSize)

	metaBytes, err := proto.Marshal(r.Meta)
	if err != nil {
		return nil, err
	}

	if metaBytes == nil {
		return nil, errMeta
	}

	rpcMetaSize := int32(len(metaBytes))
	totalSize = totalSize + rpcMetaSize

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

/*
Read byte array and initialize RpcPackage
*/
func (r *RpcDataPackage) ReadIO(rw io.Reader) error {
	if rw == nil {
		return errors.New("bytes is nil")
	}

	doInit(r)

	// read Head
	head := make([]byte, SIZE)

	_, err := io.ReadFull(rw, head)
	if err != nil {
		if err == io.EOF {
			return errIgnoreErr
		}
		log.Println("Read head error", err)
		// only to close current connection
		return err
	}

	// unmarshal Head message
	r.Head.Read(head)
	if strings.Compare(string(r.Head.MagicCode), MAGIC_CODE) != 0 {
		return fmt.Errorf("invalid magic code '%v'", string(r.Head.MagicCode))
	}

	// get RPC Meta size
	metaSize := r.Head.GetMetaSize()
	totalSize := r.Head.GetMessageSize()
	if totalSize <= 0 {
		// maybe heart beat data message, so do ignore here
		return errIgnoreErr
	}

	// read left
	leftSize := totalSize
	body := make([]byte, leftSize)

	_, err = io.ReadFull(rw, body)
	if err != nil {
		return fmt.Errorf("Read body error %w ", err)
	}

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

// Read and parse data from target byte slice
func (r *RpcDataPackage) Read(b []byte) error {

	if b == nil {
		return errors.New("parameter 'b' is nil")
	}

	buf := bytes.NewBuffer(b)

	return r.ReadIO(buf)

}

// Chunk chunk to small packages by chunk size
func (r *RpcDataPackage) Chunk(chunkSize int) []*RpcDataPackage {
	if chunkSize <= 0 {
		return []*RpcDataPackage{r}
	}

	dataSize := len(r.Data)
	chunkCount := dataSize / chunkSize
	if dataSize%chunkSize != 0 {
		chunkCount++
	}
	if chunkCount == 1 {
		return []*RpcDataPackage{r}
	}

	ret := make([]*RpcDataPackage, chunkCount)
	startPos := 0
	chunkStreamID := rand.Int63()

	for i := 0; i < chunkCount; i++ {
		temp := *r
		base := &temp // copy value

		tempMeta := *base.Meta
		base.Meta = &tempMeta
		initChuckInfo(base)
		offset := startPos + chunkSize
		if offset > dataSize {
			offset = dataSize
		}
		base.Data = base.Data[startPos:offset]
		startPos += chunkSize
		tempChuckInfo := &ChunkInfo{StreamId: chunkStreamID, ChunkId: int64(i + 1)}
		if i == chunkCount-1 {
			// this is last package
			tempChuckInfo.ChunkId = int64(-1)
		}
		if i > 0 {
			// fix duplicate attachment data
			base.Attachment = nil
		}
		ChuckInfo := *tempChuckInfo
		base.Meta.ChuckInfo = &ChuckInfo
		ret[i] = base
	}

	return ret
}

// GetChunkStreamId return chunk stream id
func (r *RpcDataPackage) GetChunkStreamId() int64 {
	initRpcMeta(r)
	return r.Meta.ChuckInfo.GetStreamId()
}

// getChunkId
func (r *RpcDataPackage) getChunkId() int64 {
	initRpcMeta(r)
	return r.Meta.ChuckInfo.GetChunkId()
}

// IsChunkPackage return if chunk package type
func (r *RpcDataPackage) IsChunkPackage() bool {
	return r.GetChunkStreamId() != 0
}

// IsFinalPackage return if the final chunk package
func (r *RpcDataPackage) IsFinalPackage() bool {
	return r.getChunkId() == -1
}

// ClearChunkStatus to clear chunk status
func (r *RpcDataPackage) ClearChunkStatus() {
	r.ChuckInfo(0, 0)
}
