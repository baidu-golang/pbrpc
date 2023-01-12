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
package baidurpc

import (
	"bytes"
	"encoding/binary"
)

const (
	SIZE = 12

	MagicSize = 4

	MAGIC_CODE = "PRPC"

	COMPRESS_NO     int32 = 0
	COMPRESS_SNAPPY int32 = 1
	COMPRESS_GZIP   int32 = 2
)

// Writable is the interface that do serialize to []byte
// if errror ocurres should return non-nil error
type Writable interface {
	Write() ([]byte, error)
}

// Readable is the interface that deserialize from []byte
// if errror ocurres should return non-nil error
type Readable interface {
	Read(bytes []byte) error
}

// RPC header content
type Header struct {
	MagicCode   []byte
	MessageSize int32
	MetaSize    int32
}

// EmptyHead return a empty head with default value
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
	bytes := make([]byte, MagicSize)
	binary.BigEndian.PutUint32(bytes, uint32(i))
	return bytes
}

// Read read byte array
func (h *Header) Read(bytes []byte) error {
	if bytes == nil || len(bytes) != SIZE {
		return nil
	}
	h.MagicCode = bytes[0:MagicSize]
	// message size offset 4 and 8
	h.MessageSize = int32(binary.BigEndian.Uint32(bytes[4:8]))
	// meta size offset 8 and 12
	h.MetaSize = int32(binary.BigEndian.Uint32(bytes[8:12]))
	return nil
}

func (h *Header) SetMagicCode(MagicCode []byte) {
	if MagicCode == nil || len(MagicCode) != MagicSize {
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
