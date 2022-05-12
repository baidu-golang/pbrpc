<!--
 * @Author: Malin Xie
 * @Description: 
 * @Date: 2021-07-24 16:54:14
-->

<h1 align="center">baidurpc</h1>

<p align="center">
baidurpc是一种基于TCP协议的二进制高性能RPC通信协议实现。它以Protobuf作为基本的数据交换格式。
本版本基于golang实现.完全兼容jprotobuf-rpc-socket: https://github.com/Baidu-ecom/Jprotobuf-rpc-socket
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/baidu-golang/pbrpc?style=flat-square)](https://goreportcard.com/report/github.com/baidu-golang/pbrpc)
[![Go](https://github.com/baidu-golang/pbrpc/actions/workflows/main.yml/badge.svg)](https://github.com/baidu-golang/pbrpc/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/baidu-golang/pbrpc/branch/master/graph/badge.svg?token=EY9Z88E82P)](https://codecov.io/gh/baidu-golang/pbrpc)
[![Releases](https://img.shields.io/github/release/baidu-golang/pbrpc/all.svg?style=flat-square)](https://github.com/baidu-golang/pbrpc/releases)
[![Go Reference](https://golang.com.cn/badge/github.com/baidu-golang/pbrpc.svg)](https://golang.com.cn/github.com/baidu-golang/pbrpc)
[![LICENSE](https://img.shields.io/github/license/baidu-golang/pbrpc.svg?style=flat-square)](https://github.com/baidu-golang/pbrpc/blob/master/LICENSE)


### features:

- 内置连接池，具备更高的性能，低延迟 QPS: 5w+
- 支持自动重连功能[Done]
- 支持附件发送[Done]
- 支持超时功能[Done]
- 压缩功能，支持GZip与Snappy[Done]
- 集成内置HTTP管理功能[TODO]
- Client支持Ha的负载均衡功能[Done]
- 灵活的超时设置功能[Done] 基于[timewheel](https://github.com/jhunters/timewheel)实现 
- 分包chunk支持，针对大数据包支持拆分包的发送的功能[Done]
- 支持Web管理能力以及内置能力[Done] [查看](https://github.com/jhunters/brpcweb)
- 支持同步发布为Http JSON协议[Done] [>= v1.2.0]
  ​
### Installing 

To start using pbrpc, install Go and run `go get`:

```sh
$ go get github.com/baidu-golang/pbrpc
```

### Which version
|version | protobuf package |
|  ----  | ----  |
|<= 1.2.x| github.com/golang/protobuf|
|1.3.x| google.golang.org/protobuf|

FYI: 由于这两个pb类库并不是完全兼容，官方推荐使用  google.golang.org/protobuf

### 使用说明与Demo 

 [Quick Start(服务发布)](./docs/quickstart_server.md) <br>
 [Quick Start(客户端调用)](./docs/quickstart_client.md) <br>
 [同步发布http rpc服务](./docs/httprpc.md) <br>
 [更多特性使用说明](./docs/Demo.md)<br>
 [Demo开发示例代码](./example)<br>
## License
brpc is [Apache 2.0 licensed](./LICENSE).
