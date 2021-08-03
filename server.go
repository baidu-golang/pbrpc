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
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/funny/link"

	"github.com/golang/protobuf/proto"
)

const (
	/** success status. */
	ST_SUCCESS int = 0

	/** 方法未找到异常. */
	ST_SERVICE_NOTFOUND int = 1001

	/** 未知异常. */
	ST_ERROR int = 2001

	/** 验证错误. */
	ST_AUTH_ERROR = 1004

	//  log id key
	KT_LOGID = "_logid_"

	RPC_STATUS_SERVICENAME = "___baidurpc_service"

	// in seconds
	Reqeust_QPS_Expire = 300 //

	Proto2_Version = "proto2"
)

// error log info definition
var (
	ERR_SERVER_NOT_INIT     = errors.New("[server-001]serverMeta is nil. please use NewTpcServer() to create TcpServer")
	ERR_INVALID_PORT        = errors.New("[server-002]port of server is nil or invalid")
	ERR_RESPONSE_TO_CLIENT  = errors.New("[server-003]response call session.Send to client failed")
	errAuth                 = errors.New("authenticate failed, pls use correct authenticate value")
	LOG_SERVICE_NOTFOUND    = "[server-" + strconv.Itoa(ST_SERVICE_NOTFOUND) + "]Service name '%s' or method name '%s' not found"
	LOG_SERVICE_DUPLICATE   = "[server-004]Service name '%s' or method name '%s' already exist"
	LOG_SERVER_STARTED_INFO = "[server-100]BaiduRpc server started on '%v'"
	LOG_INTERNAL_ERROR      = "[server-" + strconv.Itoa(ST_ERROR) + "] unknown internal error:'%s'"
	LOG_TIMECOUST_INFO      = "[server-101]Server name '%s' method '%s' process cost '%.5g' seconds"
	LOG_TIMECOUST_INFO2     = "[server-102]Server name '%s' method '%s' process cost '%.5g' seconds.(without net cost) "

	DEAFULT_IDLE_TIME_OUT_SECONDS = 10

	m           proto.Message
	MessageType = reflect.TypeOf(m)
)

type ServerMeta struct {
	Host                *string
	Port                *int
	IdleTimeoutSenconds *int
	QPSExpireInSecs     int
}

type serviceType struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

type methodType struct {
	sync.Mutex         // protects counters
	method             reflect.Method
	ArgType            reflect.Type
	ArgValue           reflect.Value
	ReturnType         reflect.Type
	ReturnValue        reflect.Value
	InArgValue         interface{}
	InPbFieldMetas     []*PbFieldMeta
	RetrunPbFieldMetas []*PbFieldMeta
}

// ParsePbMeta parse pb tag string
func (mt *methodType) ParsePbMeta() {
	t := mt.ArgType
	if t != nil {
		mt.InPbFieldMetas = parsePbMetaFromType(t)
	}

	t = mt.ReturnType
	if t != nil {
		mt.RetrunPbFieldMetas = parsePbMetaFromType(t)
	}
}

func parsePbMetaFromType(t reflect.Type) []*PbFieldMeta {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	size := t.NumField()
	if size == 0 {
		return nil
	}

	metas := make([]*PbFieldMeta, size)
	for i := 0; i < size; i++ {
		tagstruct := t.Field(i).Tag

		// check if if map
		pbtag := tagstruct.Get("protobuf")
		meta := parseMetaString(pbtag)
		if meta != nil {
			metas[i] = meta
			mapKey := tagstruct.Get("protobuf_key")
			mapValue := tagstruct.Get("protobuf_val")
			if len(mapKey) > 0 && len(mapValue) > 0 {
				meta.SubFieldMeta = make([]*PbFieldMeta, 2)
				meta.SubFieldMeta[0] = parseMetaString(mapKey)
				meta.SubFieldMeta[1] = parseMetaString(mapValue)
				meta.HasSub = true
			}

			subType := t.Field(i).Type
			if matched, _ := isMessageType(subType); matched {
				meta.SubFieldMeta = parsePbMetaFromType(subType)
			}
		}

	}
	return metas
}

func parseMetaString(meta string) *PbFieldMeta {
	partials := strings.Split(meta, ",")
	metaSize := len(partials)
	if metaSize >= 4 {
		tag, _ := strconv.Atoi(partials[1])
		meta := &PbFieldMeta{Type: partials[0], Tag: tag, Opt: partials[2]}
		// check if has version Version: partials[4]
		if metaSize >= 5 {
			meta.Version = partials[4]
		} else {
			meta.Version = Proto2_Version
		}
		nameSplit := strings.Split(partials[3], "=")
		if len(nameSplit) == 2 {
			meta.Name = nameSplit[1]
		}
		return meta
	} else {
		Warningf("invalid proto tag '%s' size is %d", meta, len(partials))
	}
	return nil
}

type PbFieldMeta struct {
	Name         string         `json:"name,omitempty"`
	Tag          int            `json:"tag,omitempty"`
	Type         string         `json:"type,omitempty"`
	Opt          string         `json:"opt,omitempty"`     // opt or req
	Version      string         `json:"version,omitempty"` // proto2 or proto3
	SubFieldMeta []*PbFieldMeta `json:"sub_field_meta,omitempty"`
	HasSub       bool           `json:"has_sub,omitempty"`
}

type attachement struct {
}

var attachementKey attachement
var errorKey struct{}

type RPCFN func(msg proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error)

// Service rpc service
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

// AuthService authenticate service
type AuthService interface {
	// Authenticate do auth action if return true auth success
	Authenticate(service, name string, authToken []byte) bool
}

// DefaultService default implemention for Service interface
type DefaultService struct {
	sname    string
	mname    string
	callback RPCFN
	inType   proto.Message
}

// DoService do call back function on RPC invocation
func (s *DefaultService) DoService(message proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error) {
	return s.callback(message, attachment, logId)
}

// GetMethodName get method name
func (s *DefaultService) GetMethodName() string {
	return s.mname
}

// NewParameter no long will be used
func (s *DefaultService) NewParameter() proto.Message {
	return s.inType
}

// GetServiceName get service name
func (s *DefaultService) GetServiceName() string {
	return s.sname
}

type Server interface {
	Start() error
	Stop() error
	Register(service *Service) (bool, error)
}

// ErrorContext
type ErrorContext struct {
	err error
}

type TcpServer struct {
	serverMeta   *ServerMeta
	services     map[string]Service
	servicesMeta map[string]*serviceMeta
	started      bool
	stop         bool
	server       *link.Server

	requestStatus *RPCRequestStatus

	authService AuthService
}

type serviceMeta struct {
	InPbFieldMetas     []*PbFieldMeta
	RetrunPbFieldMetas []*PbFieldMeta
}

// NewTpcServer
func NewTpcServer(serverMeta *ServerMeta) *TcpServer {
	tcpServer := TcpServer{}

	tcpServer.services = make(map[string]Service)
	tcpServer.servicesMeta = make(map[string]*serviceMeta)
	tcpServer.started = false
	tcpServer.stop = false

	if serverMeta.IdleTimeoutSenconds == nil {
		serverMeta.IdleTimeoutSenconds = &DEAFULT_IDLE_TIME_OUT_SECONDS
	}

	if serverMeta.QPSExpireInSecs <= 0 {
		serverMeta.QPSExpireInSecs = Reqeust_QPS_Expire
	}

	tcpServer.serverMeta = serverMeta

	// register status rpc method
	hsv := &HttpStatusView{server: &tcpServer}

	tcpServer.RegisterName(RPC_STATUS_SERVICENAME, hsv)

	return &tcpServer
}

// StartServer start server with net.Listener
func (s *TcpServer) StartServer(l net.Listener) error {

	protocol := &RpcDataPackageProtocol{}
	server := link.NewServer(l, protocol, 0 /* sync send */, link.HandlerFunc(s.handleResponse))

	s.server = server
	go server.Serve()

	s.started = true
	s.stop = false
	Infof(LOG_SERVER_STARTED_INFO, l.Addr())

	s.requestStatus = NewRPCRequestStatus(s.services) // inital request status monitor
	s.requestStatus.expireAfterSecs = int16(s.serverMeta.QPSExpireInSecs)
	go s.requestStatus.Start()
	return nil
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

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	return s.StartServer(listener)
}

func (s *TcpServer) StartAndBlock() error {
	err := s.Start()
	if err != nil {
		return err
	}
	defer s.Stop()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	fmt.Println("Press Ctrl+C or send kill sinal to exit.")
	<-c

	return nil
}

// SetAuthService set authenticate service
func (s *TcpServer) SetAuthService(authservice AuthService) {
	s.authService = authservice
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

		if s.authService != nil {
			authOk := s.authService.Authenticate(serviceName, methodName, r.Meta.AuthenticationData)
			if !authOk {
				wrapResponse(r, ST_AUTH_ERROR, errAuth.Error())
				err = session.Send(r)
				if err != nil {
					Error(ERR_RESPONSE_TO_CLIENT.Error(), "sessionId=", session.ID(), err)
				}
				return
			}
		}

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

		now2 := time.Now()
		ec := &ErrorContext{}
		// do moinitor
		s.requestStatus.RequestIn(serviceId, now2, 1)
		messageRet, attachment, err := doServiceInvoke(ec, msg, r, service)
		if ec.err != nil {
			err = ec.err
		}
		if err != nil {
			wrapResponse(r, ST_ERROR, err.Error())
			err = session.Send(r)
			if err != nil {
				Error(ERR_RESPONSE_TO_CLIENT.Error(), "sessionId=", session.ID(), err)
			}
			return
		}
		took2 := TimetookInSeconds(now2.Unix())
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

func doServiceInvoke(c *ErrorContext, msg proto.Message, r *RpcDataPackage, service Service) (proto.Message, []byte, error) {
	var err error
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("RPC server '%v' method '%v' got a internal error: %v", *r.Meta.Request.ServiceName, *r.Meta.Request.MethodName, p)
			c.err = err
			log.Println(err.Error())
		}
	}()
	messageRet, attachment, err := service.DoService(msg, r.GetAttachment(), proto.Int64(int64(r.GetLogId())))
	return messageRet, attachment, err
}

func wrapResponse(r *RpcDataPackage, errorCode int, errorText string) {
	r.GetMeta().Response = &Response{}

	r.GetMeta().GetResponse().ErrorCode = proto.Int32(int32(errorCode))
	r.GetMeta().GetResponse().ErrorText = proto.String(errorText)
}

func GetServiceId(serviceName, methodName string) string {
	return serviceName + "!" + methodName
}

func (s *TcpServer) Stop() error {
	s.stop = true
	s.started = false
	if s.server != nil {
		s.server.Stop()
	}
	if s.requestStatus != nil {
		s.requestStatus.Stop()
	}
	return nil
}

// Register register RPC service
func (s *TcpServer) Register(service interface{}) (bool, error) {
	return s.RegisterName("", service)
}

// Register register RPC service
func (s *TcpServer) registerServiceType(ss Service) (bool, error) {
	serviceId := GetServiceId(ss.GetServiceName(), ss.GetMethodName())
	exsit := s.services[serviceId]
	if exsit != nil {
		err := fmt.Errorf(LOG_SERVICE_DUPLICATE, ss.GetServiceName(), ss.GetMethodName())
		Error(err.Error())
		return false, err
	}
	log.Println("Rpc service registered. service=", ss.GetServiceName(), " method=", ss.GetMethodName())
	s.services[serviceId] = ss
	return true, nil
}

// RegisterNameWithMethodMapping call RegisterName with method name mapping map
func (s *TcpServer) RegisterNameWithMethodMapping(name string, rcvr interface{}, mapping map[string]string) (bool, error) {
	ss, ok := rcvr.(Service)
	if !ok {
		return s.registerWithMethodMapping(name, rcvr, mapping)
	}

	if name != "" {
		callback := func(msg proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error) {
			return ss.DoService(msg, attachment, logId)
		}
		mName := ss.GetMethodName()
		if mapping != nil {
			mname, ok := mapping[mName]
			if ok {
				mName = mname
			}
		}
		service := &DefaultService{
			sname:    name,
			mname:    mName,
			callback: callback,
			inType:   ss.NewParameter(),
		}
		ss = service
	}

	return s.registerServiceType(ss)
}

// RegisterName register publishes in the server with specified name for its set of methods of the
// receiver value that satisfy the following conditions:
//	- exported method of exported type
//	- one argument, exported type  and should be the type implements from proto.Message
//	- one return value, of type proto.Message
// It returns an error if the receiver is not an exported type or has
// no suitable methods. It also logs the error using package log.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.
func (s *TcpServer) RegisterName(name string, rcvr interface{}) (bool, error) {
	return s.RegisterNameWithMethodMapping(name, rcvr, nil)
}

// registerWithMethodMapping call RegisterName with method name mapping map
func (s *TcpServer) registerWithMethodMapping(name string, rcvr interface{}, mapping map[string]string) (bool, error) {
	st := &serviceType{
		typ:  reflect.TypeOf(rcvr),
		rcvr: reflect.ValueOf(rcvr),
	}

	sname := reflect.Indirect(st.rcvr).Type().Name()
	if name != "" {
		sname = name
	}
	st.name = sname

	// Install the methods
	st.method = suitableMethods(st.typ, true)

	if len(st.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PtrTo(st.typ), false)
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		log.Print(str)
		return false, errors.New(str)
	}

	// do register rpc
	for _, methodType := range st.method {
		function := methodType.method.Func
		hasArgType := false
		if methodType.ArgType != nil {
			hasArgType = true
		}
		callback := func(msg proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error) {
			// process context value
			c := context.Background()
			if attachment != nil {
				c = context.WithValue(c, attachementKey, attachment)
			}
			contextValue := reflect.ValueOf(c)

			var attachmentRet []byte = nil
			var err error

			var returnValues []reflect.Value
			if hasArgType {
				returnValues = function.Call([]reflect.Value{st.rcvr, contextValue, reflect.ValueOf(msg)})
			} else {
				returnValues = function.Call([]reflect.Value{st.rcvr, contextValue})
			}

			if len(returnValues) == 1 {
				return returnValues[0].Interface().(proto.Message), attachmentRet, nil
			}
			if len(returnValues) == 2 {
				ctx := returnValues[1].Interface().(context.Context)
				attachmentRet = Attachement(ctx)
				err = Errors(ctx)
				return returnValues[0].Interface().(proto.Message), attachmentRet, err
			}
			return nil, attachmentRet, nil
		}
		var inType proto.Message
		if methodType.InArgValue != nil {
			inType = methodType.InArgValue.(proto.Message)
			if inType == nil {
				// if not of type proto.Message
				continue
			}

		}
		mName := methodType.method.Name
		if mapping != nil {
			mname, ok := mapping[mName]
			if ok {
				mName = mname
			}
		}
		s.RegisterRpc(st.name, mName, callback, inType)

		methodType.ParsePbMeta()
		sid := GetServiceId(st.name, mName)
		s.servicesMeta[sid] = &serviceMeta{methodType.InPbFieldMetas, methodType.RetrunPbFieldMetas}
	}

	// function := mtype.method.Func
	// Invoke the method, providing a new value for the reply.
	// returnValues := function.Call([]reflect.Value{s.rcvr, argv, replyv})

	return true, nil
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)

		mtype := method.Type
		mname := method.Name
		var inArgValue interface{}
		var ok bool
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs two ins: receiver, context, *args.
		if mtype.NumIn() != 2 && mtype.NumIn() != 3 {
			if reportErr {
				log.Printf("rpc.Register: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			}
			continue
		}
		// and must be type of context.Context
		contextType := mtype.In(1)
		if !isContextType(contextType) {
			if reportErr {
				log.Printf("rpc.Register: argument type of method %q is not implements from context.Context: %q\n", mname, contextType)
			}
			continue
		}

		// and must be type of proto message
		var argType reflect.Type
		if mtype.NumIn() == 3 {
			argType = mtype.In(2)
			if ok, inArgValue = isMessageType(argType); !ok {
				if reportErr {
					log.Printf("rpc.Register: argument type of method %q is not implements from proto.Message: %q\n", mname, argType)
				}
				continue
			}
		}
		// Method needs one out.
		if mtype.NumOut() != 1 && mtype.NumOut() != 2 {
			if reportErr {
				log.Printf("rpc.Register: method %q has %d output parameters; needs one or two. \n", mname, mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		returnType := mtype.Out(0)
		if ok, _ := isMessageType(returnType); !ok {
			if reportErr {
				log.Printf("rpc.Register: return type of method %q is %q, must be implements from proto.Message\n", mname, returnType)
			}
			continue
		}
		if mtype.NumOut() == 2 {
			// The return type of the method must be error.
			returnContextType := mtype.Out(1)
			if !isContextType(returnContextType) {
				if reportErr {
					log.Printf("rpc.Register: return type of method %q is %q, must be implements from context.Context\n", mname, returnType)
				}
				continue
			}
		}

		methods[mname] = &methodType{method: method, ArgType: argType,
			ReturnType: returnType, InArgValue: inArgValue}
	}
	return methods
}

// Is this type implements from proto.Message
func isMessageType(t reflect.Type) (bool, interface{}) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// should not a interface type
	if t.Kind() == reflect.Interface {
		return false, nil
	}

	argv := reflect.New(t)
	v, ok := argv.Interface().(proto.Message)
	return ok, v
}

// Is this type implements from context.Context
func isContextType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	ok := t.Implements(reflect.TypeOf((*context.Context)(nil)).Elem())
	return ok

}

// RegisterRpc register Rpc direct
func (s *TcpServer) RegisterRpc(sname, mname string, callback RPCFN, inType proto.Message) (bool, error) {
	service := &DefaultService{
		sname:    sname,
		mname:    mname,
		callback: callback,
		inType:   inType,
	}
	return s.registerServiceType(service)
}

// Attachment utility function to get attachemnt from context
func Attachement(context context.Context) []byte {

	v := context.Value(attachementKey)
	if v == nil {
		return nil
	}
	return v.([]byte)
}

// BindAttachement add attachement value to the context
func BindAttachement(c context.Context, attachement interface{}) context.Context {
	return context.WithValue(c, attachementKey, attachement)
}

// BindError add error value to the context
func BindError(c context.Context, err error) context.Context {
	return context.WithValue(c, errorKey, err)
}

// BindError add error value to the context
func Errors(c context.Context) error {
	v := c.Value(errorKey)
	if v == nil {
		return nil
	}
	return v.(error)
}
