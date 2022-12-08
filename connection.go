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
package baidurpc

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jhunters/link"
)

var (
	errSessionIsNil  = errors.New("[conn-001]Session is nil, maybe not init Connect() function")
	errInvalidUrl    = errors.New("[conn-002]parameter 'url' of host property is nil")
	LOG_INVALID_PORT = "[conn-003]invalid parameter 'url' of port property is '%d'"
)

/*
 Connection handler interface
*/
type Connection interface {
	SendReceive(rpcDataPackage *RpcDataPackage) (*RpcDataPackage, error)
	Send(rpcDataPackage *RpcDataPackage) error
	Receive() (*RpcDataPackage, error)
	Close() error
	Reconnect() error
}

type ConnectionTester interface {
	TestConnection() error
}

// TCPConnection simple tcp based connection implementation
type TCPConnection struct {
	address      string
	sendChanSize int
	session      *link.Session[*RpcDataPackage, *RpcDataPackage]
	protocol     *RpcDataPackageProtocol[*RpcDataPackage, *RpcDataPackage]
}

/*
 Create a new TCPConnection and try to connect to target server by URL.
*/
func NewTCPConnection(url URL, timeout *time.Duration) (*TCPConnection, error) {
	connection := TCPConnection{}

	err := connection.connect(url, timeout, 0)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (c *TCPConnection) connect(url URL, timeout *time.Duration, sendChanSize int) error {
	host := url.Host
	if host == nil || len(*host) == 0 {
		return errInvalidUrl
	}
	port := url.Port
	if port == nil || *port <= 0 {
		return fmt.Errorf(fmt.Sprintf(LOG_INVALID_PORT, port))
	}

	u := *host + ":" + strconv.Itoa(*port)
	protocol, err := NewRpcDataPackageProtocol()
	if err != nil {
		return err
	}

	c.protocol = protocol
	c.protocol.timeout = timeout
	c.address = u
	c.sendChanSize = sendChanSize
	session, err := doConnect(u, protocol, timeout, sendChanSize)
	if err != nil {
		return err
	}
	c.session = session
	return nil
}

func doConnect[S, R *RpcDataPackage](address string, protocol *RpcDataPackageProtocol[*RpcDataPackage, *RpcDataPackage], timeout *time.Duration, sendChanSize int) (*link.Session[*RpcDataPackage, *RpcDataPackage], error) {
	var session *link.Session[*RpcDataPackage, *RpcDataPackage]
	var err error
	if timeout == nil {
		session, err = link.Dial[*RpcDataPackage, *RpcDataPackage]("tcp", address, protocol, sendChanSize)
	} else {
		session, err = link.DialTimeout[*RpcDataPackage, *RpcDataPackage]("tcp", address, *timeout, protocol, sendChanSize)
	}
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (c *TCPConnection) TestConnection() error {
	if c.session == nil {
		return errSessionIsNil
	}
	closed := c.session.IsClosed()
	if closed {
		return fmt.Errorf("session closed")
	}
	return nil
}

func (c *TCPConnection) GetId() uint64 {
	if c.session != nil {
		return c.session.ID()
	}

	return uint64(0)
}

// SendReceive send data to connect and block wait data recevie
func (c *TCPConnection) SendReceive(rpcDataPackage *RpcDataPackage) (*RpcDataPackage, error) {
	if c.session == nil {
		return nil, errSessionIsNil
	}

	err := c.session.Send(rpcDataPackage)

	if err != nil {
		return nil, err
	}

	return doReceive(c.session)

}

// Send data to connection
func (c *TCPConnection) Send(rpcDataPackage *RpcDataPackage) error {
	if c.session == nil {
		return errSessionIsNil
	}

	return c.session.Send(rpcDataPackage)

}

// Receive data from connection
func (c *TCPConnection) Receive() (*RpcDataPackage, error) {
	if c.session == nil {
		return nil, errSessionIsNil
	}

	return doReceive(c.session)
}

func doReceive(session *link.Session[*RpcDataPackage, *RpcDataPackage]) (rpcDataPackage *RpcDataPackage, err error) {
	rsp, err := session.Receive()
	if err != nil {
		return nil, err
	}

	if rsp == nil { // receive a error data could be ignored
		return nil, nil
	}

	return rsp, nil

}

// Close close connection
func (c *TCPConnection) Close() error {
	if c.session != nil {
		Info("close session id=", c.session.ID())
		return c.session.Close()
	}

	if c.protocol != nil {
		c.protocol.Stop()
	}
	return nil
}

// Reconnect do connect by saved info
func (c *TCPConnection) Reconnect() error {

	session, err := doConnect(c.address, c.protocol, c.protocol.timeout, c.sendChanSize)
	if err != nil {
		return err
	}
	c.session = session
	return nil
}
