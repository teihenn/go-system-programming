package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is running localhost:8080")
	for {
		// Accept()は新しいクライアント接続を受け付けるまでブロックする
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		// 各接続に対して新しいgoroutineを起動して並行処理
		go func() {
			fmt.Println("Accept %v\n", conn.RemoteAddr())
			request, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				panic(err)
			}
			dump, err := httputil.DumpRequest(request, true)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(dump))
			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Body:       io.NopCloser(strings.NewReader("Hello, World\n")),
			}
			response.Write(conn)
			conn.Close()

		}()
	}
}
