/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-11-25 11:20:51
 */
package nettool_test

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/baidu-golang/pbrpc/nettool"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	nstring string = "KHJOSJDSDKSKD"
)

func TestNewCustomListenerSelector(t *testing.T) {

	Convey("Test NewCustomListenerSelector", t, func() {
		Convey("Test NewCustomListenerSelector with host and port", func() {
			selector, err := nettool.NewCustomListenerSelector("tcp", "", 1031, 8, nettool.Equal_Mode)
			So(err, ShouldBeNil)
			So(selector, ShouldNotBeNil)
			go selector.Serve()
			defer selector.Close()
		})
		Convey("Test NewCustomListenerSelector with address string", func() {
			selector, err := nettool.NewCustomListenerSelectorByAddr("tcp", ":1032", 8, nettool.Equal_Mode)
			So(err, ShouldBeNil)
			So(selector, ShouldNotBeNil)
			go selector.Serve()
			defer selector.Close()
		})
		Convey("Test NewCustomListenerSelector with listener", func() {
			l, err := net.Listen("tcp", ":1033")
			So(err, ShouldBeNil)
			selector, err := nettool.NewCustomListenerSelectorByListener(l, 8, nettool.Equal_Mode)
			So(err, ShouldBeNil)
			So(selector, ShouldNotBeNil)
			go selector.Serve()
			defer selector.Close()
		})
	})
}

func TestEqualModeWay(t *testing.T) {
	Convey("Test Equal Mode ", t, func() {

		keystring := "XML" + nstring

		selector, err := nettool.NewCustomListenerSelector("tcp", "", 1035, 3, nettool.Equal_Mode)
		So(err, ShouldBeNil)
		So(selector, ShouldNotBeNil)

		xmlListener, err := selector.RegisterListener("XML")
		So(err, ShouldBeNil)
		So(xmlListener, ShouldNotBeNil)

		defaultLisnter := selector.RegisterDefaultListener()
		So(defaultLisnter, ShouldNotBeNil)

		go selector.Serve()
		defer selector.Close()

		// start to listen connection from socket
		defer func() {
			xmlconn, err := xmlListener.Accept()
			So(err, ShouldBeNil)
			So(xmlconn, ShouldNotBeNil)
			handleConn(xmlconn, keystring) // to process data
			xmlconn.Close()
		}()

		time.Sleep(1 * time.Second)
		// build a client connection
		conn, err := net.DialTimeout("tcp", ":1035", 2*time.Second)
		So(err, ShouldBeNil)

		// set data from client with "XML" header magic code
		_, err = conn.Write([]byte(keystring))
		So(err, ShouldBeNil)

		time.Sleep(1 * time.Second)
		conn.Close()
	})
}

func handleConn(conn net.Conn, quitString string) {
	defer conn.Close()

	result := bytes.NewBuffer(nil)
	var buf [1024]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				fmt.Println("read err:", err)
				break
			}
		} else {
			fmt.Println("recv:", result.String())
			if strings.Compare(result.String(), quitString) == 0 {
				return
			}
		}
		result.Reset()
	}
}
