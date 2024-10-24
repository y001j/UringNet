# UringNet

[![Go Report Card](https://goreportcard.com/badge/github.com/y001j/UringNet)](https://goreportcard.com/report/github.com/y001j/UringNet)
![https://codecov.io/gh/jceb/vim-orgmode](https://codecov.io/gh/y001j/uringnet/branch/master/graph/badge.svg)
[![Build](https://github.com/y001j/uringnet/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/y001j/uringnet/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/y001j/uringnet?status.svg)](https://pkg.go.dev/github.com/y001j/uringnet)

![logo](https://pic.maienzx.com/qiniuPic/202310031719836.png)

## Introduction

### Why I created this 

UringNet is a high performance network I/O framework. It is light weighted and developed by Golang. The foundation of UringNet is [io_uring](https://kernel.dk/io_uring.pdf) - a new async IO interface - from Linux Kernel since version 5.1 in 2019. This project just comes from experimental projects in my edge computing and IoT research. I am just want to find a way to build a simple but high performance network package transmit tool in IoT gateways. I have tried traditional select/epoll and then tried io_uring, finally , the io_uring get much better performance. In most cases, the performance improvement ranges from **10% to 30%**, while in the most extreme case, the performance improvement can reach up to 66%.

### Should You adopt it in your project

Though UringNet is designed to be used in IoT platforms originally, it does provide basic TCP and UDP network functions. If your project really need to handle extremely huge data traffic, you can give it a try. 

In most cases, you could just use more mature network frameworks, for example, Go [net](https://golang.org/pkg/net/)/[http](https://pkg.go.dev/net/http), [netty](https://netty.io/), [libuv](https://github.com/libuv/libuv).

## Get Started

UringNet extensively references existing network frameworks such as gnet and Netty, with some usage patterns closely resembling those of gnet. Please note that UringNet is based on io_uring and, therefore, requires running on the Linux operating system with a kernel version higher than or equal to 5.1.

### Install UringNet:

```shell
go get -u github.com/y001j/uringnet
```

### echo server example

```go
package main

import (
	"os"
	"sync"

	UringNet "github.com/y001j/uringnet"
	socket "github.com/y001j/uringnet/sockets"
)

type testServer struct {
	UringNet.BuiltinEventEngine

	testloop *UringNet.Ringloop
	//ring      *uring_net.URingNet
	addr      string
	multicore bool
}

// OnTraffic 
//  
//	@Description: OnTraffic is a hook function that runs every read event completed
//	@receiver ts
//	@param data
//	@return uring_net.Action
func (ts *testServer) OnTraffic(data *UringNet.UserData, ringnet *UringNet.URingNet) UringNet.Action {
	data.WriteBuf = data.Buffer
	return UringNet.Echo
}

func (ts *testServer) OnWritten(data UringNet.UserData) UringNet.Action {
	return UringNet.None
}

func (ts *testServer) OnOpen(data *UringNet.UserData) ([]byte, UringNet.Action) {

	return nil, UringNet.None
}

func main() {
	addr := os.Args[1]


	options := socket.SocketOptions{TCPNoDelay: socket.TCPNoDelay, ReusePort: true}
	ringNets, _ := UringNet.NewMany(UringNet.NetAddress{socket.Tcp4, addr}, 3200, true, 8, options, &testServer{}) //runtime.NumCPU()

	loop := UringNet.SetLoops(ringNets, 3000)
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	loop.RunMany()
	waitgroup.Wait()
}

```

### Important APIs

```go
UringNet.NewMany(UringNet.NetAddress{socket.Tcp4, addr}, 3200, true, 8, options, &testServer{})
```

## Benchmark

### Echo Stress Testing

In echo stress testing, the benchmark compared [go net](https://golang.org/pkg/net/), [gnet](https://github.com/panjf2000/gnet) and UringNet

Echo test tool: [rust echo bench](https://github.com/haraldh/rust_echo_bench)

> **Test environment:**
>
> Host:
>
> CPU：Intel Core i5 12400F 6 core 12 thread
>
> Memory: DDR4 32GB
>
> OS：Windows 11
>
> Virtual Machine:
>
> VMware® Workstation 16 Pro
>
> Processors: 8P ; Memory 16GB; Intel VT-x
>
> OS：Ubuntu 22.04.1 LTS， Kernel：5.19.3

#### Test Result

![echo](https://pic.maienzx.com/qiniuPic/202211061218743.png)

### HTTP Stress Testing

In HTTP stress testing, the benchmark compared [gnet](https://github.com/panjf2000/gnet), fasthttp and UringNet

Echo test tool: [wrk](https://github.com/wg/wrk)

> Test Enviroment 
>
> Host：
>
> CPU：Intel Core i5 12400F 6 core 12 thread
>
> Memory: DDR4 32GB
>
> OS：Windows 11
>
> Testing Virtual Machine：
>
> VMware® Workstation 16 Pro
>
> Processors: 8P ; Memory 16GB; Intel VT-x
>
> OS：Ubuntu 22.04.1 LTS， Kernel：5.19.3
>
> wrk Virtual Machine：
>
> VMware® Workstation 16 Pro
>
> Processors: 4P ; Memory 8GB; Intel VT-x
>
> OS：Ubuntu 20.04.5 LTS， Kernel：5.15.0   

#### Test Result

![image-20221212113551004](https://pic.maienzx.com/qiniuPic/image-20221212113551004.png)

## Known Issues

> **Warning:** The project is still in development, there are still many bugs and performance issues. If you find any bugs or have any suggestions, please feel free to open an issue.

During our performance tests with wrk, we observed that UringNet requires a warm-up period during HTTP stress testing. Typically, a warm-up time of around 5 minutes is necessary to ensure that the UringNet test code reaches its peak performance, which surpasses the current fastest framework by more than 10%. I am still investigating the root cause of this behavior.
