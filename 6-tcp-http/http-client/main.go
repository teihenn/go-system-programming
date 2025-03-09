package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}
	request.Write(conn)
	// ネットワークからのレスポンスデータはまず OS の TCP バッファに蓄積される。
	// bufio.NewReader を使うことで、TCP バッファから 4KB ずつまとめて読み込み、
	// conn.Read() のシステムコール回数を削減し、パフォーマンスを向上させる。
	response, err := http.ReadResponse(bufio.NewReader(conn), request)
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
}
