/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-19 13:22:01
 */
package baidurpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	HttpRpcPath = "/rpc/"

	LogId_key        = "X-LogID"
	Auth_key         = "X-Authenticate"
	Trace_Id_key     = "X-Trace_ID"
	Trace_Span_key   = "X-Trace_Span"
	Trace_Parent_key = "X-Trace_Parent"
	Request_Meta_Key = "X-Request-Meta" // Json value
)

// ResponseData
type ResponseData struct {
	ErrNo   int         `json:"errno"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HttpServer struct {
	s       *TcpServer
	httpsrv *http.Server
}

func (h *HttpServer) serverHttp(l net.Listener) {

	srv := &http.Server{
		Handler: h,
	}

	h.httpsrv = srv

	go func() {
		// service connections
		if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

}

// ServeHTTP to serve http reqeust and response to process http rpc handle
func (h *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if !strings.HasPrefix(path, HttpRpcPath) {
		data := toJson(errResponse(ST_SERVICE_NOTFOUND, fmt.Sprintf("no service or method found by path='%s'", path)))
		w.Write(data)
		return
	}
	serviceName, method, err := getServiceMethod(path)
	if err != nil {
		data := toJson(errResponse(ST_ERROR, err.Error()))
		w.Write(data)
		return
	}

	sid := GetServiceId(serviceName, method)

	service, ok := h.s.services[sid]
	if !ok {
		data := toJson(errResponse(ST_SERVICE_NOTFOUND, fmt.Sprintf("no service or method found by path='%s'", path)))
		w.Write(data)
		return
	}

	// authenticate
	if h.s.authService != nil {
		authData := getHeaderAsByte(req, Auth_key)
		authOk := h.s.authService.Authenticate(serviceName, method, authData)
		if !authOk {
			data := toJson(errResponse(ST_AUTH_ERROR, errAuth.Error()))
			w.Write(data)
			return
		}
	}

	if h.s.traceService != nil {
		traceId := getHeaderAsInt64(req, Trace_Id_key)
		spanId := getHeaderAsInt64(req, Trace_Span_key)
		parentId := getHeaderAsInt64(req, Trace_Parent_key)
		traceInfo := &TraceInfo{TraceId: traceId, SpanId: spanId, ParentSpanId: parentId}

		value := getHeaderAsByte(req, Request_Meta_Key)
		if value != nil {
			value, _ = UnescapeUnicode(value)
			metaExt := map[string]string{}
			err := json.Unmarshal(value, &metaExt)
			if err == nil {
				traceInfo.RpcRequestMetaExt = metaExt
			}
		}

		traceRetrun := h.s.traceService.Trace(serviceName, method, traceInfo)
		if traceRetrun != nil {
			w.Header().Set(Trace_Id_key, int64ToString(traceRetrun.TraceId))
			w.Header().Set(Trace_Span_key, int64ToString(traceRetrun.SpanId))
			w.Header().Set(Trace_Parent_key, int64ToString(traceRetrun.ParentSpanId))
			if traceRetrun.RpcRequestMetaExt != nil {
				metaData, err := json.Marshal(traceRetrun.RpcRequestMetaExt)
				if err == nil {
					w.Header().Set(Request_Meta_Key, string(metaData))
				}
			}
		}
	}

	// get json data
	jsonData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		data := toJson(errResponse(ST_ERROR, err.Error()))
		w.Write(data)
		return
	}

	paramIn := service.NewParameter()
	err = json.Unmarshal(jsonData, paramIn)
	if err != nil {
		data := toJson(errResponse(ST_ERROR, err.Error()))
		w.Write(data)
		return
	}

	// get logid
	var logid int64 = getHeaderAsInt64(req, LogId_key)

	ec := &ErrorContext{}
	ret, _, err := h.s.doServiceInvoke(ec, paramIn, serviceName, method, nil, logid, service)
	if err != nil {
		data := toJson(errResponse(ST_ERROR, err.Error()))
		w.Write(data)
		return
	}

	resData := &ResponseData{ErrNo: 0, Data: ret}
	data, err := json.Marshal(resData)
	if err != nil {
		data := toJson(errResponse(ST_ERROR, err.Error()))
		w.Write(data)
		return
	}

	w.Write(data)

}

// shutdown do shutdown action to close http server
func (h *HttpServer) shutdown(ctx context.Context) {
	if h.httpsrv != nil {
		h.httpsrv.Shutdown(ctx)
	}
}

func getHeaderAsByte(req *http.Request, key string) []byte {
	value, ok := req.Header[http.CanonicalHeaderKey(key)]
	if !ok || len(value) != 1 {
		return nil
	}
	return []byte(value[0])
}

func getHeaderAsInt64(req *http.Request, key string) int64 {
	value, ok := req.Header[http.CanonicalHeaderKey(key)]
	if !ok || len(value) != 1 {
		return -1
	}

	id, _ := strconv.Atoi(value[0])
	return int64(id)
}

func int64ToString(i int64) string {
	return strconv.Itoa(int(i))
}

func errResponse(errno int, message string) *ResponseData {
	return &ResponseData{ErrNo: errno, Message: message}
}

func toJson(o interface{}) []byte {
	data, _ := json.Marshal(o)
	return data
}

func getServiceMethod(path string) (string, string, error) {
	p := strings.TrimPrefix(path, HttpRpcPath)
	p = strings.TrimSuffix(p, "/")

	seperate := strings.Split(p, "/")
	if len(seperate) != 2 {
		return "", "", fmt.Errorf("no service or method found by path='%s'", p)
	}
	return seperate[0], seperate[1], nil
}

// UnescapeUnicode
func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
