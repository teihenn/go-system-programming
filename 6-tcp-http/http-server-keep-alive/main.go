package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// main関数：HTTP/1.1のKeep-Alive機能を使用したHTTPサーバーを実装
// Keep-Aliveとは、一度確立したTCP接続を複数のHTTPリクエスト/レスポンスで再利用する機能。
// 接続の確立と切断のオーバーヘッドを減らし、パフォーマンスを向上させる。
func main() {
	// TCPリスナーを作成し、localhost:8080で待ち受ける
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		// エラーが発生した場合はプログラムを終了
		panic(err)
	}
	fmt.Println("Server is running localhost:8080")
	for {
		// 新しい接続を受け入れる
		conn, err := listener.Accept()
		if err != nil {
			// エラーが発生した場合はプログラムを終了
			panic(err)
		}
		// 新しいゴルーチンを開始して接続を処理
		go func() {
			defer conn.Close() // 関数終了時に接続を閉じる
			fmt.Println("Accept %v\n", conn.RemoteAddr())
			// TCPコネクションが張られたあとに何度もリクエストを受けられるようにforで回す
			for {
				// 読み込みのタイムアウトを5秒に設定
				// これはKeep-Alive接続が無限に開いたままにならないようにするため
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				// HTTPリクエストを読み込む
				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					// タイムアウトまたはソケットクローズ時はループを終了。
					// その他のエラーはパニックを発生させる。
					// error型から、より具体的なnet.Error型に変換する。
					// net.Errorインターフェースには、Timeout()メソッドが定義されており、タイムアウトエラーかどうかを判断できる。
					// ok変数はキャストが成功したかどうかを示す（成功=true、失敗=false）
					neterr, ok := err.(net.Error) // エラーをネットワークエラーとしてダウンキャスト
					if ok && neterr.Timeout() {
						fmt.Println("Timeout")
						break
					} else if err == io.EOF {
						// io.EOFはEnd Of Fileの略で、データの終わりを示すエラー。
						// HTTP通信では、クライアントが接続を明示的に閉じた場合に発生する。
						// Keep-Alive接続において、クライアントが正常に接続を終了したことを検出し、
						// サーバー側もループを抜けて接続処理を終了する。
						break
					}
					panic(err)
				}
				// リクエスト内容をダンプして表示
				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(dump))
				content := "Hello World\n"

				// HTTPレスポンスを作成
				// HTTP/1.1であることと、ContentLengthの設定が必要
				// ContentLengthはレスポンスの終わりをクライアントに伝えるために必須
				// これが無いとクライアントはレスポンスの終了を判断できず、Keep-Aliveが正常に機能しない
				response := http.Response{
					StatusCode:    200,                                      // ステータスコード200（OK）
					ProtoMajor:    1,                                        // HTTP/1.xのメジャーバージョン
					ProtoMinor:    1,                                        // HTTP/1.xのマイナーバージョン（HTTP/1.1でKeep-Aliveがデフォルトとなる）
					ContentLength: int64(len(content)),                      // コンテンツの長さ
					Body:          io.NopCloser(strings.NewReader(content)), // レスポンスボディ
				}
				// レスポンスを接続に書き込む
				response.Write(conn)
			}
		}()
	}
}
