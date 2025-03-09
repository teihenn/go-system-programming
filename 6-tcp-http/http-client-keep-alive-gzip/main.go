package main

import (
	"bufio"
	"compress/gzip" // gzip圧縮データを解凍するためのパッケージ
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func main() {
	sendMessages := []string{
		"ASCII",
		"PROGRAMMING",
		"PLUS",
	}
	current := 0
	var conn net.Conn = nil
	for {
		var err error
		if conn == nil {
			conn, err = net.Dial("tcp", "localhost:8080")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Access: %d\n", current)
		}
		request, err := http.NewRequest(
			"POST",
			"http://localhost:8080",
			strings.NewReader(sendMessages[current]))
		if err != nil {
			panic(err)
		}
		// サーバーにgzip圧縮したレスポンスを要求するヘッダーを設定
		// gzip圧縮を使うことで：
		// - 転送データ量が削減され、帯域幅を効率的に使用できる
		// - 特に大きなテキストデータでは転送速度が向上する
		// - モバイル通信など帯域制限のある環境でも効率的に通信できる
		request.Header.Set("Accept-Encoding", "gzip")
		err = request.Write(conn)
		if err != nil {
			panic(err)
		}
		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			fmt.Println("Retry")
			conn = nil
			continue
		}
		// レスポンスヘッダーのみ表示（bodyはfalseに設定）
		dump, err := httputil.DumpResponse(response, false)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		// レスポンスボディを確実に閉じる。
		// これをしないとメモリリーク、ファイルディスクリプタの枯渇、
		// コネクションプール不足などの深刻な問題が発生する。
		// ボディは実際のデータストリームを表すI/Oリソースであり、
		// ヘッダーと違って明示的に解放が必要
		defer response.Body.Close()

		// gzip圧縮されたレスポンスかどうかを確認して適切に処理
		// サーバーが実際にgzip圧縮を適用したかはContent-Encodingヘッダーで判断する
		// gzip圧縮されたデータは特殊なバイナリ形式なので、そのまま使用せず解凍が必要
		if response.Header.Get("Content-Encoding") == "gzip" {
			// gzip圧縮されている場合は解凍用のリーダーを作成
			// gzipリーダーは元のデータを自動的に解凍しながら読み取る変換ストリーム
			reader, err := gzip.NewReader(response.Body)
			if err != nil {
				panic(err)
			}
			// 解凍したデータを標準出力に書き込む
			io.Copy(os.Stdout, reader)
			// gzipリーダーを明示的に閉じる
			// これによりリソースリークを防止し、内部バッファもクリーンアップされる
			reader.Close()
		} else {
			// 圧縮されていない場合はそのまま標準出力に書き込む
			io.Copy(os.Stdout, response.Body)
		}

		current++
		if current == len(sendMessages) {
			break
		}
	}
	conn.Close()
}
