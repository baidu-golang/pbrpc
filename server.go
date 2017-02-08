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
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/baidu-golang/pbrpc/link"
	"github.com/golang/protobuf/proto"
)

const (
	/** success status. */
	ST_SUCCESS int = 0

	/** 方法未找到异常. */
	ST_SERVICE_NOTFOUND int = 1001

	/** 未知异常. */
	ST_ERROR int = 2001
)

// error log info definition
var ERR_SERVER_NOT_INIT = errors.New("[server-001]serverMeta is nil. please use NewTpcServer() to create TcpServer.")
var ERR_INVALID_PORT = errors.New("[server-002]port of server is nil or invalid.")
var ERR_RESPONSE_TO_CLIENT = errors.New("[server-003]response call session.Send to client failed.")
var LOG_SERVICE_NOTFOUND = "[server-" + strconv.Itoa(ST_SERVICE_NOTFOUND) + "]Service name '%s' or method name '%s' not found."
var LOG_SERVICE_DUPLICATE = "[server-004]Service name '%s' or method name '%s' already exist."
var LOG_SERVER_STARTED_INFO = "[server-100]Server started on '%s' of port '%d'"
var LOG_INTERNAL_ERROR = "[server-" + strconv.Itoa(ST_ERROR) + "] unknown internal error:'%s'"
var LOG_TIMECOUST_INFO = "[server-101]Server name '%s' method '%s' process cost '%.5g' seconds."
var LOG_TIMECOUST_INFO2 = "[server-102]Server name '%s' method '%s' process cost '%.5g' seconds.(without net cost) "

var DEAFULT_IDLE_TIME_OUT_SECONDS = 10

type ServerMeta struct {
	Host                *string
	Port                *int
	IdleTimeoutSenconds *int
}

type Service interface {
	/*
	   RPC service call back method.
	   message : parameter in from RPC client or 'nil' if has no parameter
	   attachment : attachment content from RPC client or 'nil' if has no attachment
	   logId : with a int64 type log sequence id from client or 'nil if has no logId
	   return:
	   [0] message return back to RPC client or 'nil' if need not return method response
	   [1] attachment return back to RPC client or 'nil' if need not return attachemnt
	   [2] return with any error or 'nil' represents success
	*/
	DoService(message proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error)
	GetServiceName() string
	GetMethodName() string
	NewParameter() proto.Message
}

type Server interface {
	Start() error
	Stop() error
	Register(service *Service) (bool, error)
}

type TcpServer struct {
	serverMeta *ServerMeta
	services   map[string]Service
	started    bool
	stop       bool
	server     *link.Server
}

func NewTpcServer(serverMeta *ServerMeta) *TcpServer {
	tcpServer := TcpServer{}

	tcpServer.services = make(map[string]Service)
	tcpServer.started = false
	tcpServer.stop = false

	if serverMeta.IdleTimeoutSenconds == nil {
		serverMeta.IdleTimeoutSenconds = &DEAFULT_IDLE_TIME_OUT_SECONDS
	}

	tcpServer.serverMeta = serverMeta

	return &tcpServer
}

func (s *TcpServer) Start() error {
	if s.serverMeta == nil {
		return ERR_SERVER_NOT_INIT
	}

	var addr = ""
	host := ""
	if s.serverMeta.Host != nil {
		host = *s.serverMeta.Host
	}

	port := s.serverMeta.Port
	if port == nil || *port <= 0 {
		return ERR_INVALID_PORT
	}

	addr = addr + host + ":" + strconv.Itoa(*port)

	protocol := &RpcDataPackageProtocol{}
	server, err := link.Listen("tcp", addr, protocol, 0 /* sync send */, link.HandlerFunc(s.handleResponse))

	if err != nil {
		return err
	}
	server.IdleTimeoutSenconds = s.serverMeta.IdleTimeoutSenconds
	s.server = server
	go server.Serve()

	s.started = true
	defer func() {
		s.started = false
	}()

	Infof(LOG_SERVER_STARTED_INFO, host, *port)

	return nil
}

func (s *TcpServer) StartAndBlock() error {
	err := s.Start()
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	fmt.Println("Press Ctrl+C or send kill sinal to exit.")
	_ = <-c

	return nil
}

func (s *TcpServer) handleResponse(session *link.Session) {
	// after function return must close session
	defer session.Close()

	for {

		now := time.Now().UnixNano()

		req, err := session.Receive()
		if err != nil {
			Errorf(LOG_INTERNAL_ERROR, err.Error())
			return
		}
		// error package
		if req == nil {
			return
		}

		r, ok := req.(*RpcDataPackage)
		if !ok {
			return // convert error maybe type mismatch
		}

		serviceName := r.GetMeta().GetRequest().GetServiceName()
		methodName := r.GetMeta().GetRequest().GetMethodName()

		serviceId := GetServiceId(serviceName, methodName)

		service := s.services[serviceId]
		if service == nil {
			wrapResponse(r, ST_SERVICE_NOTFOUND, fmt.Sprintf(LOG_SERVICE_NOTFOUND, serviceName, methodName))

			err = session.Send(r)
			if err != nil {
				Error(ERR_RESPONSE_TO_CLIENT.Error(), "sessionId=", session.ID(), err)
			}
			return
		}

		var msg proto.Message
		requestData := r.GetData()
		if requestData != nil {
			msg = service.NewParameter()
			if msg != nil {
				proto.Unmarshal(requestData, msg)
			}
		}
		// do service here
		now2 := time.Now().Unix()
		messageRet, attachment, err := service.DoService(msg, r.GetAttachment(), proto.Int64(int64(r.GetLogId())))
		if err != nil {
			wrapResponse(r, ST_ERROR, err.Error())
			err = session.Send(r)
			if err != nil {
				Error(ERR_RESPONSE_TO_CLIENT.Error(), "sessionId=", session.ID(), err)
			}
			return
		}
		took2 := TimetookInSeconds(now2)
		Infof(LOG_TIMECOUST_INFO2, serviceName, methodName, took2)

		if messageRet == nil {
			r.SetData(nil)
		} else {
			d, err := proto.Marshal(messageRet)
			if err != nil {
				wrapResponse(r, ST_ERROR, err.Error())
				err = session.Send(r)
				if err != nil {
					Error(ERR_RESPONSE_TO_CLIENT.Error(), "sessionId=", session.ID(), err)
				}
				return
			}

			r.SetData(d)
			r.SetAttachment(attachment)
			wrapResponse(r, ST_SUCCESS, "")
		}

		err = session.Send(r)

		if err != nil {
			Error(ERR_RESPONSE_TO_CLIENT.Error(), "sessionId=", session.ID(), err)
			return
		}

		took := TimetookInSeconds(now)
		Infof(LOG_TIMECOUST_INFO, serviceName, methodName, took)

	}

}

func wrapResponse(r *RpcDataPackage, errorCode int, errorText string) {
	r.GetMeta().Response = &Response{}

	r.GetMeta().GetResponse().ErrorCode = proto.Int(errorCode)
	r.GetMeta().GetResponse().ErrorText = proto.String(errorText)
}

func GetServiceId(serviceName, methodName string) string {
	return serviceName + "!" + methodName
}

func (s *TcpServer) Stop() error {
	s.stop = true
	if s.server != nil {
		s.server.Stop()
	}
	return nil
}

func (s *TcpServer) Register(service Service) (bool, error) {
	ss := service
	serviceId := GetServiceId(ss.GetServiceName(), ss.GetMethodName())
	exsit := s.services[serviceId]
	if exsit != nil {
		err := errors.New(fmt.Sprintf(LOG_SERVICE_DUPLICATE, ss.GetServiceName(), ss.GetMethodName()))
		Error(err.Error())
		return false, err
	}
	s.services[serviceId] = ss
	return true, nil
}
