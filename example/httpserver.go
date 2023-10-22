package main

import (
	"bytes"
	"github.com/y001j/UringNet"
	"github.com/y001j/UringNet/sockets"
	"os"
	"sync"
	"time"
)

type testServer struct {
	UringNet.BuiltinEventEngine

	testloop *UringNet.Ringloop
	//ring      *uring_net.URingNet
	addr      string
	multicore bool
}

type httpCodec struct {
	delimiter []byte
	buf       []byte
}

func appendResponse(hc *[]byte) {
	*hc = append(*hc, "HTTP/1.1 200 OK\r\nServer: uringNet\r\nContent-Type: text/plain\r\nDate: "...)
	*hc = time.Now().AppendFormat(*hc, "Mon, 02 Jan 2006 15:04:05 GMT")
	*hc = append(*hc, "\r\nContent-Length: 12\r\n\r\nHello World!"...)
}

var (
	errMsg      = "Internal Server Error"
	errMsgBytes = []byte(errMsg)
)

func (ts *testServer) OnTraffic(data *UringNet.UserData, ringnet UringNet.URingNet) UringNet.Action {

	//将data.Buffer转换为string
	//buffer := data.Buffer[:data.BufSize]

	buffer := ringnet.ReadBuffer
	//tes :=
	//fmt.Println("data:", " offset: ", tes, " ", data.BufOffset)
	//获取tes中“\r\n\r\n”的数量
	count := bytes.Count(buffer, []byte("GET"))
	if count == 0 {
		//appendResponse(&data.WriteBuf)
		//return UringNet.Close
		return UringNet.None
	} else {
		for i := 0; i < count; i++ {
			appendResponse(&data.WriteBuf)
		}
	}
	return UringNet.Echo
}

func (ts *testServer) OnWritten(data UringNet.UserData) UringNet.Action {

	return UringNet.None
}

//func (hc *httpCodec) parse(data []byte) (int, error) {
//	if idx := bytes.Index(data, hc.delimiter); idx != -1 {
//		return idx + 4, nil
//	}
//	return -1, errCRLFNotFound
//}
//
//func (hc *httpCodec) reset() {
//	hc.buf = hc.buf[:0]
//}
//
//var errCRLFNotFound = errors.New("CRLF not found")

func (ts *testServer) OnOpen(data *UringNet.UserData) ([]byte, UringNet.Action) {

	ts.SetContext(&httpCodec{delimiter: []byte("\r\n\r\n")})
	return nil, UringNet.None
}

func main() {
	addr := os.Args[1]
	//runtime.GOMAXPROCS(runtime.NumCPU()*2 - 1)

	options := socket.SocketOptions{TCPNoDelay: socket.TCPNoDelay, ReusePort: true}
	ringNets, _ := UringNet.NewMany(UringNet.NetAddress{socket.Tcp4, addr}, 3200, true, 4, options, &testServer{}) //runtime.NumCPU()

	loop := UringNet.SetLoops(ringNets, 4000)

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)

	loop.RunMany()

	waitgroup.Wait()
}
