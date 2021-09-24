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
/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-07-23 18:46:22
 */
package nettool

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	Equal_Mode     = 1
	StartWith_Mode = 2
)

type NetInfo struct {
	conn net.Conn
	err  error
}

// CustomListener wrap Listener interface
type CustomListener struct {
	listenerProxy net.Listener

	sessionChan chan NetInfo
}

func (l *CustomListener) Accept() (net.Conn, error) {
	conn := <-l.sessionChan
	return conn.conn, conn.err
}

func (l *CustomListener) Close() error {
	return l.listenerProxy.Close()
}

func (l *CustomListener) Addr() net.Addr {
	return l.listenerProxy.Addr()
}

// CustomListenerSelector a selector CustomListener by head value
type CustomListenerSelector struct {
	headSize        uint8
	defaultListener *CustomListener
	listeners       map[string]*CustomListener

	listenerProxy net.Listener
	matchMode     int
}

// NewCustomListenerSelector new a CustomListenerSelector
func NewCustomListenerSelector(network, host string, port int, headsize uint8, matchMode int) (*CustomListenerSelector, error) {
	return NewCustomListenerSelectorByAddr(network, host+":"+strconv.Itoa(port), headsize, matchMode)
}

// NewCustomListenerSelectorByAddr new a CustomListenerSelector by address string
func NewCustomListenerSelectorByAddr(network, server string, headsize uint8, matchMode int) (*CustomListenerSelector, error) {
	listener, err := net.Listen(network, server)
	if err != nil {
		return nil, err
	}

	return NewCustomListenerSelectorByListener(listener, headsize, matchMode)
}

// NewCustomListenerSelectorByListener
func NewCustomListenerSelectorByListener(l net.Listener, headsize uint8, matchMode int) (*CustomListenerSelector, error) {
	if matchMode != Equal_Mode && matchMode != StartWith_Mode {
		return nil, fmt.Errorf("invalid match mode value %d", matchMode)
	}

	selector := &CustomListenerSelector{listenerProxy: l, headSize: headsize, matchMode: matchMode}
	selector.listeners = make(map[string]*CustomListener)

	selector.defaultListener = &CustomListener{listenerProxy: l, sessionChan: make(chan NetInfo)}

	return selector, nil
}

// RegisterListener
func (server *CustomListenerSelector) RegisterListener(headMagiccode string) (net.Listener, error) {
	if len(headMagiccode) != int(server.headSize) && server.matchMode == Equal_Mode {
		return nil, fmt.Errorf("error head magic code '%s', size should be '%d'", headMagiccode, server.headSize)
	}
	listener := &CustomListener{listenerProxy: server.listenerProxy, sessionChan: make(chan NetInfo)}
	server.listeners[headMagiccode] = listener
	return listener, nil
}

// RegisterDefaultListener
func (server *CustomListenerSelector) RegisterDefaultListener() net.Listener {
	return server.defaultListener
}

// Serve do listening from net trasport
func (server *CustomListenerSelector) Serve() error {
	for {
		conn, err := server.listenerProxy.Accept()
		if err != nil {
			log.Println("CustomListenerSelector started failed.", err)
			// if met error broadcast to all listeners
			netinfo := NetInfo{conn, err}
			for _, server := range server.listeners {
				s := server
				go func() {
					s.sessionChan <- netinfo
				}()
			}
			return err
		}

		go func() {
			r := conn.(io.Reader)

			// read Head
			head := make([]byte, server.headSize)
			r.Read(head)
			cw := &ConnWrapper{conn, head, 0, server.headSize}
			netinfo := NetInfo{cw, nil}

			if server.matchMode == Equal_Mode {
				serverListener, ok := server.listeners[string(head)]
				if !ok { // not exist, use default listener
					server.defaultListener.sessionChan <- netinfo
				} else {
					serverListener.sessionChan <- netinfo
				}
			} else { // start with mode
				var matched bool = false
				for magicCode, listener := range server.listeners {
					if strings.HasPrefix(string(head), magicCode) {
						listener.sessionChan <- netinfo
						matched = true
						break
					}
				}
				if !matched { // not matched use default listener
					server.defaultListener.sessionChan <- netinfo
				}
			}

		}()

	}

}

// Close do close all listeners
func (server *CustomListenerSelector) Close() error {
	var errRet error
	for _, server := range server.listeners {
		err := server.Close()
		if err != nil {
			errRet = err
		}
	}

	return errRet
}

//  net.Conn  proxy
type ConnWrapper struct {
	conn     net.Conn
	head     []byte
	n        int
	headsize uint8
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (cw *ConnWrapper) Read(b []byte) (n int, err error) {
	size := len(b)

	if size > int(cw.headsize) && cw.n < int(cw.headsize) {
		for i := cw.n; i < int(cw.headsize); i++ {
			b[i] = cw.head[i]
		}
		cw.n = size
		nn, err := cw.conn.Read(b[cw.headsize:])
		if err != nil {
			return nn + int(cw.headsize), err
		}
		return nn + int(cw.headsize), nil
	} else if size <= int(cw.headsize) && cw.n < int(cw.headsize) {
		for i := cw.n; i < size; i++ {
			b[i] = cw.head[i]
		}
		cw.n = size
		return size, nil
	}
	// otherwize read directly
	count, err := cw.conn.Read(b)
	return count, err
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (cw *ConnWrapper) Write(b []byte) (n int, err error) {
	return cw.conn.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (cw *ConnWrapper) Close() error {
	return cw.conn.Close()
}

// LocalAddr returns the local network address.
func (cw *ConnWrapper) LocalAddr() net.Addr {
	return cw.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (cw *ConnWrapper) RemoteAddr() net.Addr {
	return cw.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (cw *ConnWrapper) SetDeadline(t time.Time) error {
	return cw.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (cw *ConnWrapper) SetReadDeadline(t time.Time) error {
	return cw.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (cw *ConnWrapper) SetWriteDeadline(t time.Time) error {
	return cw.conn.SetWriteDeadline(t)
}
