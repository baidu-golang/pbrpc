/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-07-26 17:09:25
 */
package baidurpc

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jhunters/timewheel"
)

// HttpStatusView
type HttpStatusView struct {
	server *TcpServer
}

type RPCStatus struct {
	Host            string       `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port            int32        `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	TimeoutSenconds int32        `protobuf:"varint,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Methods         []*RPCMethod `protobuf:"bytes,4,rep,name=methods,proto3" json:"methods,omitempty"`
}

type RPCMethod struct {
	Service string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Method  string `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
}

func (m *RPCStatus) Reset()         { *m = RPCStatus{} }
func (m *RPCStatus) String() string { return proto.CompactTextString(m) }
func (*RPCStatus) ProtoMessage()    {}

func (m *RPCMethod) Reset()         { *m = RPCMethod{} }
func (m *RPCMethod) String() string { return proto.CompactTextString(m) }
func (*RPCMethod) ProtoMessage()    {}

func (hsv *HttpStatusView) Status(c context.Context) (*RPCStatus, context.Context) {
	s := hsv.server
	result := &RPCStatus{}
	if s.serverMeta.Host != nil {
		result.Host = *s.serverMeta.Host
	}
	if s.serverMeta.Port != nil {
		result.Port = int32(*s.serverMeta.Port)
	}
	if s.serverMeta.IdleTimeoutSenconds != nil {
		result.TimeoutSenconds = int32(*s.serverMeta.IdleTimeoutSenconds)
	}

	rpcServices := s.services
	methods := make([]*RPCMethod, len(rpcServices))
	var i int = 0
	for _, service := range rpcServices {
		m := &RPCMethod{Service: service.GetServiceName(), Method: service.GetMethodName()}
		methods[i] = m
		i++
	}
	result.Methods = methods
	return result, c
}

// RPCRequestStatus
type RPCRequestStatus struct {
	Methods map[string]*RPCMethodReuqestStatus

	reqeustChan chan request

	closeChan chan bool

	expireAfterSecs int16

	started bool

	tw *timewheel.TimeWheel
}

type request struct {
	method string
	t      time.Time
}

// RPCMethodReuqestStatus
type RPCMethodReuqestStatus struct {
	QpsStatus map[int64]int32
}

// NewRPCRequestStatus
func NewRPCRequestStatus(services map[string]Service) *RPCRequestStatus {
	ret := &RPCRequestStatus{
		Methods:     make(map[string]*RPCMethodReuqestStatus, len(services)),
		reqeustChan: make(chan request, 1024),
		closeChan:   make(chan bool),
	}

	for name := range services {
		ret.Methods[name] = &RPCMethodReuqestStatus{QpsStatus: make(map[int64]int32, 1024)}
	}

	return ret
}

// Start
func (r *RPCRequestStatus) Start() error {
	Info("RPC method eeuqest status record starting")
	r.started = true

	// start time wheel to delete expire data
	tw, err := timewheel.New(1*time.Second, uint16(r.expireAfterSecs))
	r.tw = tw
	r.tw.Start()
	if err != nil {
		r.started = false
		return err
	}

	for {
		select {
		case m := <-r.reqeustChan:
			status, ok := r.Methods[m.method]
			if !ok {
				status = &RPCMethodReuqestStatus{QpsStatus: make(map[int64]int32, 1024)}
				r.Methods[m.method] = status
			}
			k := m.t.Unix()
			count, ok := status.QpsStatus[k]
			if !ok {
				count = 1
			} else {
				count++
			}
			status.QpsStatus[k] = count

			fmt.Println("request record", m.method, count)
		case <-r.closeChan:
			r.started = false
			return nil
		}
	}

}

// RequestIn
func (r *RPCRequestStatus) RequestIn(methodName string, t time.Time) error {
	if !r.started {
		return fmt.Errorf("RequestIn failed status not started")
	}
	req := request{method: methodName, t: t}
	r.reqeustChan <- req

	task := timewheel.Task{
		Data: req,
		TimeoutCallback: func(tt timewheel.Task) { // call back function on time out
			k := tt.Data.(request)
			r.expire(k.method, k.t)

		}}

	// add task and return unique task id
	_, err := r.tw.AddTask(time.Duration(r.expireAfterSecs)*time.Second, task) // add delay task
	return err
}

func (r *RPCRequestStatus) expire(methodName string, t time.Time) {
	status, ok := r.Methods[methodName]
	if ok {
		delete(status.QpsStatus, t.Unix())
	}
}

// Close
func (r *RPCRequestStatus) Close() {
	if !r.started {
		return
	}
	r.started = false
	r.closeChan <- true

	r.tw.Stop()
}
