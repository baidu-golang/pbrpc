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
		data := toJson(errResponse(fmt.Sprintf("no service or method found by path='%s'", path)))
		w.Write(data)
		return
	}

	serviceName, method, err := getServiceMethod(path)
	if err != nil {
		data := toJson(errResponse(err.Error()))
		w.Write(data)
		return
	}

	sid := GetServiceId(serviceName, method)

	service, ok := h.s.services[sid]
	if !ok {
		data := toJson(errResponse(fmt.Sprintf("no service or method found by path='%s'", path)))
		w.Write(data)
		return
	}

	// get json data
	jsonData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		data := toJson(errResponse(err.Error()))
		w.Write(data)
		return
	}

	paramIn := service.NewParameter()
	err = json.Unmarshal(jsonData, paramIn)
	if err != nil {
		data := toJson(errResponse(err.Error()))
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

	ret, _, err := service.DoService(paramIn, nil, logid)
	if err != nil {
		data := toJson(errResponse(err.Error()))
		w.Write(data)
		return
	}

	resData := &ResponseData{ErrNo: 0, Data: ret}
	data, err := json.Marshal(resData)
	if err != nil {
		data := toJson(errResponse(err.Error()))
		w.Write(data)
		return
	}

	w.Write(data)

}

func errResponse(message string) *ResponseData {
	return &ResponseData{ErrNo: -1, Message: message}
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
