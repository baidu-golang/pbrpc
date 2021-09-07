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
	"context"
	"errors"
	"time"

	pool "github.com/jolestar/go-commons-pool/v2"
)

var (
	Empty_Head = make([]byte, SIZE)

	errPoolNotInit      = errors.New("[connpool-001]Object pool is nil maybe not init Connect() function")
	errGetConnFail      = errors.New("[connpool-002]Can not get connection from connection pool. target object is nil")
	errDestroyObjectNil = errors.New("[connpool-003]Destroy object failed due to target object is nil")

	HB_SERVICE_NAME = "__heartbeat"
	HB_METHOD_NAME  = "__beat"
)

/*
type Connection interface {
	SendReceive(rpcDataPackage *RpcDataPackage) (*RpcDataPackage, error)
	Close() error
}
*/

type TCPConnectionPool struct {
	Config     *pool.ObjectPoolConfig
	objectPool *pool.ObjectPool
	timeout    *time.Duration
}

func NewDefaultTCPConnectionPool(url URL, timeout *time.Duration) (*TCPConnectionPool, error) {
	return NewTCPConnectionPool(url, timeout, nil)
}

func NewTCPConnectionPool(url URL, timeout *time.Duration, config *pool.ObjectPoolConfig) (*TCPConnectionPool, error) {
	connection := TCPConnectionPool{timeout: timeout}
	if config == nil {
		connection.Config = pool.NewDefaultPoolConfig()
		connection.Config.TestOnBorrow = true

	} else {
		connection.Config = config
	}

	err := connection.connect(url, timeout, 0)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (c *TCPConnectionPool) connect(url URL, timeout *time.Duration, sendChanSize int) error {

	factory := ConnectionPoolFactory{timeout: timeout}
	factory.url = &url

	var objectPool *pool.ObjectPool

	eConfig := c.Config
	if eConfig == nil {
		eConfig = pool.NewDefaultPoolConfig()
	}

	objectPool = pool.NewObjectPool(context.Background(), &factory, eConfig)
	c.objectPool = objectPool

	return nil

}

func (c *TCPConnectionPool) borrowObject() (*TCPConnection, error) {
	if c.objectPool == nil {
		return nil, errPoolNotInit
	}

	object, err := c.objectPool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}

	if object == nil {
		return nil, errGetConnFail
	}

	return object.(*TCPConnection), nil
}

func (c *TCPConnectionPool) SendReceive(rpcDataPackage *RpcDataPackage) (*RpcDataPackage, error) {
	object, err := c.borrowObject()
	if err != nil {
		return nil, err
	}
	defer c.objectPool.ReturnObject(context.Background(), object)

	return object.SendReceive(rpcDataPackage)

}

// Send
func (c *TCPConnectionPool) Send(rpcDataPackage *RpcDataPackage) error {
	object, err := c.borrowObject()
	if err != nil {
		return err
	}
	defer c.objectPool.ReturnObject(context.Background(), object)
	err = object.Send(rpcDataPackage)
	return err

}

// Receive
func (c *TCPConnectionPool) Receive() (*RpcDataPackage, error) {
	object, err := c.borrowObject()
	if err != nil {
		return nil, err
	}
	defer c.objectPool.ReturnObject(context.Background(), object)

	return object.Receive()
}

func (c *TCPConnectionPool) Close() error {
	if c.objectPool == nil {
		return errPoolNotInit
	}

	c.objectPool.Close(context.Background())
	return nil
}

func (c *TCPConnectionPool) GetNumActive() int {
	if c.objectPool == nil {
		return 0
	}

	return c.objectPool.GetNumActive()
}

// Reconnect do connect by saved info
func (c *TCPConnectionPool) Reconnect() error {
	// do nothing
	return nil
}

type ConnectionPoolFactory struct {
	url     *URL
	timeout *time.Duration
}

func (c *ConnectionPoolFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	if c.url == nil {
		return nil, errPoolNotInit
	}

	connection := TCPConnection{}
	err := connection.connect(*c.url, c.timeout, 0)
	if err != nil {
		return nil, err
	}

	return pool.NewPooledObject(&connection), nil
}

func (c *ConnectionPoolFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	obj := object.Object
	if obj == nil {
		return errDestroyObjectNil
	}

	conn := obj.(Connection)
	return conn.Close()
}

func (c *ConnectionPoolFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {

	obj := object.Object
	if obj == nil {
		return false
	}

	conn := obj.(ConnectionTester)

	return conn.TestConnection() == nil
}

func (c *ConnectionPoolFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (c *ConnectionPoolFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
