package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {
	// 送信するメッセージのリスト
	sendMessages := []string{
		"ASCII",
		"PROGRAMMING",
		"PLUS",
	}
	// 現在処理中のメッセージインデックス
	current := 0
	// TCP接続を保持する変数。複数リクエストで再利用するためにループの外で定義
	var conn net.Conn = nil
	// リトライ用にループで全体を囲う
	for {
		var err error
		// まだコネクションを張っていない／エラーでリトライする場合の処理
		// keep-aliveの要：接続が存在しない場合のみ新規に接続する
		if conn == nil {
			conn, err = net.Dial("tcp", "localhost:8080")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Access: %d\n", current)
		}
		// POSTで文字列を送るリクエストを作成
		request, err := http.NewRequest(
			"POST",
			"http://localhost:8080",
			strings.NewReader(sendMessages[current]))
		if err != nil {
			panic(err)
		}
		// 既存の接続にリクエストを書き込む
		err = request.Write(conn)
		if err != nil {
			panic(err)
		}
		// サーバーから読み込む
		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			// エラー発生時は接続をリセットしてリトライする
			// これがkeep-alive対応の重要な部分：接続エラー時の再接続処理
			fmt.Println("Retry")
			conn = nil
			continue
		}
		// レスポンスの内容を表示
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))
		// 全部送信完了したら終了
		current++
		if current == len(sendMessages) {
			break
		}
		// ループを継続：同じ接続を使って次のリクエストを送信する
		// 非keep-alive版では毎回新しい接続を作成する
	}
	// 使用済みの接続を閉じる
	conn.Close()
}
