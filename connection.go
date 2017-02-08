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
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/baidu-golang/pbrpc/link"
)

var ERR_SESSION_IS_NIL = errors.New("[conn-001]Session is nil, maybe not init Connect() function.")
var ERR_INVALID_URL = errors.New("[conn-002]parameter 'url' of host property is nil")

var LOG_INVALID_PORT = "[conn-003]invalid parameter 'url' of port property is '%d'"

/*
 Connection handler interface
*/
type Connection interface {
	SendReceive(rpcDataPackage *RpcDataPackage) (*RpcDataPackage, error)
	Close() error
}

type ConnectionTester interface {
	TestConnection() error
}

type TCPConnection struct {
	session *link.Session
}

/*
 Create a new TCPConnection and try to connect to target server by URL.
*/
func NewTCPConnection(url URL, timeout *time.Duration) (*TCPConnection, error) {
	connection := TCPConnection{}

	var err error
	err = connection.connect(url, timeout, 0)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (c *TCPConnection) connect(url URL, timeout *time.Duration, sendChanSize int) error {
	host := url.Host
	if host == nil || len(*host) == 0 {
		return ERR_INVALID_URL
	}
	port := url.Port
	if port == nil || *port <= 0 {
		return errors.New(fmt.Sprintf(LOG_INVALID_PORT, port))
	}

	u := *host + ":" + strconv.Itoa(*port)
	protocol := &RpcDataPackageProtocol{}
	var session *link.Session
	var err error
	if timeout == nil {
		session, err = link.Dial("tcp", u, protocol, sendChanSize)
	} else {
		session, err = link.DialTimeout("tcp", u, *timeout, protocol, sendChanSize)
	}
	if err != nil {
		return err
	}
	c.session = session
	return nil
}

func (c *TCPConnection) TestConnection() error {
	if c.session == nil {
		return ERR_SESSION_IS_NIL
	}

	b, _ := EmptyHead().Write()
	err := c.session.Send(b)
	if err != nil {
		return err
	}
	_, err = c.session.Receive()
	return err
}

func (c *TCPConnection) GetId() uint64 {
	if c.session != nil {
		return c.session.ID()
	}

	return uint64(0)
}

func (c *TCPConnection) SendReceive(rpcDataPackage *RpcDataPackage) (*RpcDataPackage, error) {
	if c.session == nil {
		return nil, ERR_SESSION_IS_NIL
	}

	err := c.session.Send(rpcDataPackage)

	if err != nil {
		return nil, err
	}

	return doReceive(c.session)

}

func doReceive(session *link.Session) (rpcDataPackage *RpcDataPackage, err error) {
	rsp, err := session.Receive()
	if err != nil {
		return nil, err
	}

	if rsp == nil { // receive a error data could be ignored
		return nil, nil
	}

	r := rsp.(*RpcDataPackage)
	return r, nil

}

func (c *TCPConnection) Close() error {
	if c.session != nil {
		Info("close session id=", c.session.ID())
		return c.session.Close()
	}
	return nil
}
