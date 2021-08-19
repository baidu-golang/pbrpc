/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-19 13:22:01
 */
package baidurpc

import (
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

	LogId_key = "X-LogID"
)

// ResponseData
type ResponseData struct {
	ErrNo   int         `json:"errno"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HttpServer struct {
	s *TcpServer
}

func (h *HttpServer) serverHttp(l net.Listener) {

	srv := &http.Server{
		Handler: h,
	}

	go func() {
		// service connections
		if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

}

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
	var logid *int64
	sLogId, ok := req.Header[http.CanonicalHeaderKey(LogId_key)]
	if ok && len(sLogId) == 1 {
		id, _ := strconv.Atoi(sLogId[0])
		newId := int64(id)
		logid = &newId
	}

	ec := &ErrorContext{}
	ret, _, err := h.s.doServiceInvoke(ec, paramIn, serviceName, method, nil, *logid, service)
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
