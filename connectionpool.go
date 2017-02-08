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
	"time"

	pool "github.com/jolestar/go-commons-pool"
)

var Empty_Head = make([]byte, SIZE)

var ERR_POOL_NOT_INIT = errors.New("[connpool-001]Object pool is nil maybe not init Connect() function.")
var ERR_GET_CONN_FAIL = errors.New("[connpool-002]Can not get connection from connection pool. target object is nil.")
var ERR_DESTORY_OBJECT_NIL = errors.New("[connpool-003]Destroy object failed due to target object is nil.")
var ERR_POOL_URL_NIL = errors.New("[connpool-004]Can not create object cause of URL is nil.")

var HB_SERVICE_NAME = "__heartbeat"
var HB_METHOD_NAME = "__beat"

/*
type Connection interface {
	SendReceive(rpcDataPackage *RpcDataPackage) (*RpcDataPackage, error)
	Close() error
}
*/

type TCPConnectionPool struct {
	Config     *pool.ObjectPoolConfig
	objectPool *pool.ObjectPool
}

func NewDefaultTCPConnectionPool(url URL, timeout *time.Duration) (*TCPConnectionPool, error) {
	return NewTCPConnectionPool(url, timeout, nil)
}

func NewTCPConnectionPool(url URL, timeout *time.Duration, config *pool.ObjectPoolConfig) (*TCPConnectionPool, error) {
	connection := TCPConnectionPool{}
	if config == nil {
		connection.Config = pool.NewDefaultPoolConfig()
		connection.Config.TestOnBorrow = true

	} else {
		connection.Config = config
	}

	var err error
	err = connection.connect(url, timeout, 0)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (c *TCPConnectionPool) connect(url URL, timeout *time.Duration, sendChanSize int) error {

	factory := ConnectionPoolFactory{}
	factory.url = &url

	var objectPool *pool.ObjectPool

	eConfig := c.Config
	if eConfig == nil {
		eConfig = pool.NewDefaultPoolConfig()
	}

	objectPool = pool.NewObjectPool(&factory, eConfig)
	c.objectPool = objectPool

	return nil

}

func (c *TCPConnectionPool) borrowObject() (*TCPConnection, error) {
	if c.objectPool == nil {
		return nil, ERR_POOL_NOT_INIT
	}

	object, err := c.objectPool.BorrowObject()
	if err != nil {
		return nil, err
	}

	if object == nil {
		return nil, ERR_GET_CONN_FAIL
	}

	return object.(*TCPConnection), nil
}

func (c *TCPConnectionPool) SendReceive(rpcDataPackage *RpcDataPackage) (*RpcDataPackage, error) {
	object, err := c.borrowObject()
	if err != nil {
		return nil, err
	}
	defer c.objectPool.ReturnObject(object)

	return object.SendReceive(rpcDataPackage)

}

func (c *TCPConnectionPool) Close() error {
	if c.objectPool == nil {
		return ERR_POOL_NOT_INIT
	}

	c.objectPool.Close()
	return nil
}

func (c *TCPConnectionPool) GetNumActive() int {
	if c.objectPool == nil {
		return 0
	}

	return c.objectPool.GetNumActive()
}

type ConnectionPoolFactory struct {
	url *URL
}

func (c *ConnectionPoolFactory) MakeObject() (*pool.PooledObject, error) {
	if c.url == nil {
		return nil, ERR_POOL_URL_NIL
	}

	connection := TCPConnection{}
	err := connection.connect(*c.url, nil, 0)
	if err != nil {
		return nil, err
	}

	return pool.NewPooledObject(&connection), nil
}

func (c *ConnectionPoolFactory) DestroyObject(object *pool.PooledObject) error {
	obj := object.Object
	if obj == nil {
		return ERR_DESTORY_OBJECT_NIL
	}

	conn := obj.(Connection)
	return conn.Close()
}

func (c *ConnectionPoolFactory) ValidateObject(object *pool.PooledObject) bool {

	if true {
		return true
	}

	obj := object.Object
	if obj == nil {
		return false
	}

	conn := obj.(ConnectionTester)

	return conn.TestConnection() != nil
}

func (c *ConnectionPoolFactory) ActivateObject(object *pool.PooledObject) error {
	return nil
}

func (c *ConnectionPoolFactory) PassivateObject(object *pool.PooledObject) error {
	return nil
}
