// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.9.2
// source: brpc_meta.proto

package baidurpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RpcMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Request            *Request   `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	Response           *Response  `protobuf:"bytes,2,opt,name=response,proto3" json:"response,omitempty"`
	CompressType       int32      `protobuf:"varint,3,opt,name=compress_type,json=compressType,proto3" json:"compress_type,omitempty"` // 0:nocompress 1:Snappy 2:gzip
	CorrelationId      int64      `protobuf:"varint,4,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	AttachmentSize     int32      `protobuf:"varint,5,opt,name=attachment_size,json=attachmentSize,proto3" json:"attachment_size,omitempty"`
	ChuckInfo          *ChunkInfo `protobuf:"bytes,6,opt,name=chuck_info,json=chuckInfo,proto3" json:"chuck_info,omitempty"`
	AuthenticationData []byte     `protobuf:"bytes,7,opt,name=authentication_data,json=authenticationData,proto3" json:"authentication_data,omitempty"`
}

func (x *RpcMeta) Reset() {
	*x = RpcMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_brpc_meta_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RpcMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RpcMeta) ProtoMessage() {}

func (x *RpcMeta) ProtoReflect() protoreflect.Message {
	mi := &file_brpc_meta_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RpcMeta.ProtoReflect.Descriptor instead.
func (*RpcMeta) Descriptor() ([]byte, []int) {
	return file_brpc_meta_proto_rawDescGZIP(), []int{0}
}

func (x *RpcMeta) GetRequest() *Request {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *RpcMeta) GetResponse() *Response {
	if x != nil {
		return x.Response
	}
	return nil
}

func (x *RpcMeta) GetCompressType() int32 {
	if x != nil {
		return x.CompressType
	}
	return 0
}

func (x *RpcMeta) GetCorrelationId() int64 {
	if x != nil {
		return x.CorrelationId
	}
	return 0
}

func (x *RpcMeta) GetAttachmentSize() int32 {
	if x != nil {
		return x.AttachmentSize
	}
	return 0
}

func (x *RpcMeta) GetChuckInfo() *ChunkInfo {
	if x != nil {
		return x.ChuckInfo
	}
	return nil
}

func (x *RpcMeta) GetAuthenticationData() []byte {
	if x != nil {
		return x.AuthenticationData
	}
	return nil
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName       string                    `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	MethodName        string                    `protobuf:"bytes,2,opt,name=method_name,json=methodName,proto3" json:"method_name,omitempty"`
	LogId             int64                     `protobuf:"varint,3,opt,name=log_id,json=logId,proto3" json:"log_id,omitempty"`
	TraceId           int64                     `protobuf:"varint,4,opt,name=traceId,proto3" json:"traceId,omitempty"`
	SpanId            int64                     `protobuf:"varint,5,opt,name=spanId,proto3" json:"spanId,omitempty"`
	ParentSpanId      int64                     `protobuf:"varint,6,opt,name=parentSpanId,proto3" json:"parentSpanId,omitempty"`
	RpcRequestMetaExt []*RpcRequestMetaExtField `protobuf:"bytes,7,rep,name=rpcRequestMetaExt,proto3" json:"rpcRequestMetaExt,omitempty"`
	ExtraParam        []byte                    `protobuf:"bytes,110,opt,name=extraParam,proto3" json:"extraParam,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_brpc_meta_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_brpc_meta_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_brpc_meta_proto_rawDescGZIP(), []int{1}
}

func (x *Request) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *Request) GetMethodName() string {
	if x != nil {
		return x.MethodName
	}
	return ""
}

func (x *Request) GetLogId() int64 {
	if x != nil {
		return x.LogId
	}
	return 0
}

func (x *Request) GetTraceId() int64 {
	if x != nil {
		return x.TraceId
	}
	return 0
}

func (x *Request) GetSpanId() int64 {
	if x != nil {
		return x.SpanId
	}
	return 0
}

func (x *Request) GetParentSpanId() int64 {
	if x != nil {
		return x.ParentSpanId
	}
	return 0
}

func (x *Request) GetRpcRequestMetaExt() []*RpcRequestMetaExtField {
	if x != nil {
		return x.RpcRequestMetaExt
	}
	return nil
}

func (x *Request) GetExtraParam() []byte {
	if x != nil {
		return x.ExtraParam
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorCode int32  `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3" json:"error_code,omitempty"`
	ErrorText string `protobuf:"bytes,2,opt,name=error_text,json=errorText,proto3" json:"error_text,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_brpc_meta_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_brpc_meta_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_brpc_meta_proto_rawDescGZIP(), []int{2}
}

func (x *Response) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *Response) GetErrorText() string {
	if x != nil {
		return x.ErrorText
	}
	return ""
}

type ChunkInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreamId int64 `protobuf:"varint,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
	ChunkId  int64 `protobuf:"varint,2,opt,name=chunk_id,json=chunkId,proto3" json:"chunk_id,omitempty"`
}

func (x *ChunkInfo) Reset() {
	*x = ChunkInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_brpc_meta_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChunkInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChunkInfo) ProtoMessage() {}

func (x *ChunkInfo) ProtoReflect() protoreflect.Message {
	mi := &file_brpc_meta_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChunkInfo.ProtoReflect.Descriptor instead.
func (*ChunkInfo) Descriptor() ([]byte, []int) {
	return file_brpc_meta_proto_rawDescGZIP(), []int{3}
}

func (x *ChunkInfo) GetStreamId() int64 {
	if x != nil {
		return x.StreamId
	}
	return 0
}

func (x *ChunkInfo) GetChunkId() int64 {
	if x != nil {
		return x.ChunkId
	}
	return 0
}

type RpcRequestMetaExtField struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *RpcRequestMetaExtField) Reset() {
	*x = RpcRequestMetaExtField{}
	if protoimpl.UnsafeEnabled {
		mi := &file_brpc_meta_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RpcRequestMetaExtField) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RpcRequestMetaExtField) ProtoMessage() {}

func (x *RpcRequestMetaExtField) ProtoReflect() protoreflect.Message {
	mi := &file_brpc_meta_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RpcRequestMetaExtField.ProtoReflect.Descriptor instead.
func (*RpcRequestMetaExtField) Descriptor() ([]byte, []int) {
	return file_brpc_meta_proto_rawDescGZIP(), []int{4}
}

func (x *RpcRequestMetaExtField) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *RpcRequestMetaExtField) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_brpc_meta_proto protoreflect.FileDescriptor

var file_brpc_meta_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x62, 0x72, 0x70, 0x63, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xa5, 0x02, 0x0a, 0x07, 0x52, 0x70, 0x63, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x22, 0x0a,
	0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x25, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0c, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x25, 0x0a,
	0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x61,
	0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x29, 0x0a,
	0x0a, 0x63, 0x68, 0x75, 0x63, 0x6b, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0a, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x63,
	0x68, 0x75, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2f, 0x0a, 0x13, 0x61, 0x75, 0x74, 0x68,
	0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x12, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x22, 0xa1, 0x02, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x6c, 0x6f, 0x67,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x49, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x70,
	0x61, 0x6e, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x70, 0x61, 0x6e,
	0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x70, 0x61, 0x6e,
	0x49, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74,
	0x53, 0x70, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x45, 0x0a, 0x11, 0x72, 0x70, 0x63, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x78, 0x74, 0x18, 0x07, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x52, 0x70, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x65,
	0x74, 0x61, 0x45, 0x78, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x11, 0x72, 0x70, 0x63, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x78, 0x74, 0x12, 0x1e, 0x0a,
	0x0a, 0x65, 0x78, 0x74, 0x72, 0x61, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x6e, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x72, 0x61, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x22, 0x48, 0x0a,
	0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x54, 0x65, 0x78, 0x74, 0x22, 0x43, 0x0a, 0x09, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49,
	0x64, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x49, 0x64, 0x22, 0x40, 0x0a, 0x16,
	0x52, 0x70, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x78,
	0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x22,
	0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x61, 0x69,
	0x64, 0x75, 0x2d, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x62, 0x61, 0x69, 0x64, 0x75, 0x72,
	0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_brpc_meta_proto_rawDescOnce sync.Once
	file_brpc_meta_proto_rawDescData = file_brpc_meta_proto_rawDesc
)

func file_brpc_meta_proto_rawDescGZIP() []byte {
	file_brpc_meta_proto_rawDescOnce.Do(func() {
		file_brpc_meta_proto_rawDescData = protoimpl.X.CompressGZIP(file_brpc_meta_proto_rawDescData)
	})
	return file_brpc_meta_proto_rawDescData
}

var file_brpc_meta_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_brpc_meta_proto_goTypes = []interface{}{
	(*RpcMeta)(nil),                // 0: RpcMeta
	(*Request)(nil),                // 1: Request
	(*Response)(nil),               // 2: Response
	(*ChunkInfo)(nil),              // 3: ChunkInfo
	(*RpcRequestMetaExtField)(nil), // 4: RpcRequestMetaExtField
}
var file_brpc_meta_proto_depIdxs = []int32{
	1, // 0: RpcMeta.request:type_name -> Request
	2, // 1: RpcMeta.response:type_name -> Response
	3, // 2: RpcMeta.chuck_info:type_name -> ChunkInfo
	4, // 3: Request.rpcRequestMetaExt:type_name -> RpcRequestMetaExtField
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_brpc_meta_proto_init() }
func file_brpc_meta_proto_init() {
	if File_brpc_meta_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_brpc_meta_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RpcMeta); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_brpc_meta_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_brpc_meta_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_brpc_meta_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChunkInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_brpc_meta_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RpcRequestMetaExtField); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_brpc_meta_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_brpc_meta_proto_goTypes,
		DependencyIndexes: file_brpc_meta_proto_depIdxs,
		MessageInfos:      file_brpc_meta_proto_msgTypes,
	}.Build()
	File_brpc_meta_proto = out.File
	file_brpc_meta_proto_rawDesc = nil
	file_brpc_meta_proto_goTypes = nil
	file_brpc_meta_proto_depIdxs = nil
}
