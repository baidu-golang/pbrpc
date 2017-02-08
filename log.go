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
	"flag"
	"log"

	"github.com/golang/glog"
)

var gloged = flag.Bool("glog", false, "If non-empty, use default log otherwise use glog")

// to enable glog enabled
// enabled: true to enable glog otherwise use log module instead
func EnableGLog(enabled bool) {
	gloged = &enabled
}

// Info logs to the INFO log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Info(args ...interface{}) {
	if *gloged {
		glog.Info(args...)
	} else {
		log.Println(args...)
	}
}

// Infof logs to the INFO log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Infof(format string, args ...interface{}) {
	if *gloged {
		glog.Infof(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// Warning logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Warning(args ...interface{}) {
	if *gloged {
		glog.Warning(args...)
	} else {
		log.Println(args...)
	}
}

// Warningf logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	if *gloged {
		glog.Warningf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// Error logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Error(args ...interface{}) {
	if *gloged {
		glog.Error(args...)
	} else {
		log.Println(args...)
	}
}

// Errorf logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {
	if *gloged {
		glog.Errorf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
