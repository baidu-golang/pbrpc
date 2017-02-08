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
)

const SIZE = 12

const MAGIC_CODE = "PRPC"

const COMPRESS_NO int32 = 0
const COMPRESS_SNAPPY int32 = 1
const COMPRESS_GZIP int32 = 2

type Writable interface {
	Write() ([]byte, error)
}

type Readable interface {
	Read(bytes []byte) error
}

// RPC header content
type Header struct {
	MagicCode   []byte
	MessageSize int32
	MetaSize    int32
}

func EmptyHead() *Header {
	h := Header{}
	h.MagicCode = []byte(MAGIC_CODE)
	h.MessageSize = 0
	h.MetaSize = 0
	return &h
}

/*
	Convert Header struct to byte array
*/
func (h *Header) Write() ([]byte, error) {
	b := new(bytes.Buffer)

	binary.Write(b, binary.BigEndian, h.GetMagicCode())
	binary.Write(b, binary.BigEndian, intToBytes(h.GetMessageSize()))
	binary.Write(b, binary.BigEndian, intToBytes(h.GetMetaSize()))
	return b.Bytes(), nil
}

func intToBytes(i int32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(i))
	return bytes
}

func (h *Header) Read(bytes []byte) error {
	if bytes == nil || len(bytes) != SIZE {
		return nil
	}
	h.MagicCode = bytes[0:4]
	h.MessageSize = int32(binary.BigEndian.Uint32(bytes[4:8]))
	h.MetaSize = int32(binary.BigEndian.Uint32(bytes[8:12]))
	return nil
}

func (h *Header) SetMagicCode(MagicCode []byte) {
	if MagicCode == nil || len(MagicCode) != 4 {
		return
	}
	h.MagicCode = MagicCode
}

func (h *Header) GetMagicCode() []byte {
	return h.MagicCode
}

func (h *Header) SetMessageSize(MessageSize int32) {
	h.MessageSize = MessageSize
}

func (h *Header) GetMessageSize() int32 {
	return h.MessageSize
}

func (h *Header) SetMetaSize(MetaSize int32) {
	h.MetaSize = MetaSize
}

func (h *Header) GetMetaSize() int32 {
	return h.MetaSize
}
