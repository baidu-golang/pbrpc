package baidurpc

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

// HaRpcClient high avialbe RpcClient
type HaRpcClient struct {
	rpcClients []*RpcClient

	current int32

	locker sync.Mutex
}

// NewBatchTCPConnection to create batch connection
func NewBatchTCPConnection(urls []URL, timeout time.Duration) ([]Connection, error) {
	if len(urls) == 0 {
		return nil, errors.New("param 'urls' is empty")
	}

	var result []Connection
	for _, url := range urls {
		conn, err := NewTCPConnection(url, &timeout)
		if err != nil {
			log.Println("create connection failed:", err)
			continue
		}
		result = append(result, conn)
	}

	return result, nil
}

// CloseBatchConnection close batch connections
func CloseBatchConnection(connections []Connection) {
	for _, conn := range connections {
		conn.Close()
	}
}

// NewHaRpcCient
func NewHaRpcCient(connections []Connection) (*HaRpcClient, error) {
	if len(connections) == 0 {
		return nil, errors.New("param 'connections' is empty")
	}

	rpcClients := make([]*RpcClient, len(connections))
	for idx, connection := range connections {
		rpcClient, err := NewRpcCient(connection)
		if err != nil {
			return nil, err
		}
		rpcClients[idx] = rpcClient
	}
	result := &HaRpcClient{rpcClients: rpcClients, current: 0}
	return result, nil
}

// electClient
func (c *HaRpcClient) electClient() *RpcClient {
	if len(c.rpcClients) == 0 {
		return nil
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	client := c.rpcClients[int(c.current)%len(c.rpcClients)]
	c.current++
	return client
}

// SendRpcRequest
func (c *HaRpcClient) SendRpcRequest(rpcInvocation *RpcInvocation, responseMessage proto.Message) (*RpcDataPackage, error) {
	var errRet error
	size := len(c.rpcClients)
	if size == 0 {
		return nil, errors.New("no rpc client avaible")
	}

	for i := 0; i < size; i++ {
		rpcClient := c.electClient()
		if rpcClient == nil {
			return nil, errors.New("no rpc client avaible")
		}

		data, err := rpcClient.SendRpcRequest(rpcInvocation, responseMessage)
		if err == nil {
			return data, nil
		} else {
			errRet = err
		}
	}

	return nil, errRet
}
