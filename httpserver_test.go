/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-08-20 17:16:06
 */
package baidurpc_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	baidurpc "github.com/baidu-golang/pbrpc"
	. "github.com/smartystreets/goconvey/convey"
)

// TestHttpRpcServerWithPath
func TestHttpRpcServerWithPath(t *testing.T) {
	Convey("test http mode with service not found", t, func() {

		tcpServer := startRpcServerWithHttpMode(0, true)
		defer stopRpcServer(tcpServer)

		Convey("test http mode with bad prefix path", func() {
			sport := strconv.Itoa(PORT_1)
			urlpath := "http://localhost:" + sport + "/badpath/EchoService/echo"
			paramters := map[string]string{}
			headers := map[string]string{}

			result := sendHttpRequest(urlpath, "POST", "{\"name\":\"matt\"}", paramters, headers)
			So(result, ShouldNotBeNil)
			data := &baidurpc.ResponseData{}
			err := json.Unmarshal(result, data)
			So(err, ShouldBeNil)
			So(data.ErrNo, ShouldEqual, baidurpc.ST_SERVICE_NOTFOUND)
		})

		Convey("test http mode with bad service and method path", func() {
			sport := strconv.Itoa(PORT_1)
			urlpath := "http://localhost:" + sport + "/rpc/BadService/BadMethod"
			paramters := map[string]string{}
			headers := map[string]string{}

			result := sendHttpRequest(urlpath, "POST", "{\"name\":\"matt\"}", paramters, headers)
			So(result, ShouldNotBeNil)
			data := &baidurpc.ResponseData{}
			err := json.Unmarshal(result, data)
			So(err, ShouldBeNil)
			So(data.ErrNo, ShouldEqual, baidurpc.ST_SERVICE_NOTFOUND)
		})

	})
}

// TestHttpRpcServer
func TestHttpRpcServer(t *testing.T) {
	Convey("test http rpc with common request", t, func() {

		tcpServer := startRpcServerWithHttpMode(0, true)
		defer stopRpcServer(tcpServer)

		sport := strconv.Itoa(PORT_1)
		urlpath := "http://localhost:" + sport + "/rpc/EchoService/echo"
		paramters := map[string]string{}
		headers := map[string]string{}

		result := sendHttpRequest(urlpath, "POST", "{\"name\":\"matt\"}", paramters, headers)

		So(result, ShouldNotBeNil)
		data := &baidurpc.ResponseData{}
		err := json.Unmarshal(result, data)
		So(err, ShouldBeNil)
		So(data.ErrNo, ShouldEqual, 0)

	})
}

// TestHttpRpcServerWithAuthenticate
func TestHttpRpcServerWithAuthenticate(t *testing.T) {
	Convey("test http rpc with authenticate", t, func() {

		tcpServer := startRpcServerWithHttpMode(0, true)
		tcpServer.SetAuthService(new(StringMatchAuthService))
		tcpServer.SetTraceService(new(AddOneTraceService))
		defer stopRpcServer(tcpServer)

		sport := strconv.Itoa(PORT_1)
		urlpath := "http://localhost:" + sport + "/rpc/EchoService/echo"
		paramters := map[string]string{}
		headers := map[string]string{}
		headers[baidurpc.Auth_key] = AUTH_TOKEN
		headers[baidurpc.LogId_key] = "1"

		headers[baidurpc.Trace_Id_key] = "1"
		headers[baidurpc.Trace_Span_key] = "2"
		headers[baidurpc.Trace_Parent_key] = "3"

		result := sendHttpRequest(urlpath, "POST", "{\"name\":\"matt\"}", paramters, headers)

		So(result, ShouldNotBeNil)
		data := &baidurpc.ResponseData{}
		err := json.Unmarshal(result, data)
		So(err, ShouldBeNil)
		So(data.ErrNo, ShouldEqual, 0)

	})
}

func sendHttpRequest(urlpath, method, body string, paramters, headers map[string]string) []byte {
	params := url.Values{}
	Url, err := url.Parse(urlpath)
	if err != nil {
		return nil
	}

	for k, v := range paramters {
		params.Set(k, v)
	}

	urlPath := Url.String()

	client := &http.Client{}
	req, _ := http.NewRequest(method, urlPath, nil)

	data := []byte(body)
	bs := bytes.NewBuffer(data)
	req.ContentLength = int64(len(data))
	// buf := body.Bytes()
	req.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(data)
		return ioutil.NopCloser(r), nil
	}

	req.Body = ioutil.NopCloser(bs)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, _ := client.Do(req)
	rbody, _ := ioutil.ReadAll(resp.Body)

	return rbody
}
