package main

import (
	"net"
	"fmt"
	"bufio"
)

func main () {
	// Echoサーバーを立ち上げる
	// まず、Listenで8080番ポートを使うことを予約する
	// Listenの戻り値はlistenインターフェース？
	listener, err := net.Listen("tcp", ":8080") // ":8080"の：って何？
	if err != nil {
		fmt.Println("起動失敗！")
		return
	}
	// 後で8080番ポートを閉じておく
	defer listener.Close()

	// Acceptはfor文の中、外？
	for {
		// Acceptの引数は？
		// クライアントの接続を受け付ける
		// Acceptの戻り値はなんなの？
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接続失敗")
			continue
			// 接続に失敗した後はreturnしちゃダメじゃない？
			// 他の接続も全部切ってプログラム終了するのはおかしい
		}
		// connって何？
		// スレッドを作り、裏で関数を動かす
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// サーバー側の処理だよね？
	// 引数がconnなのはなぜ？
	scanner := bufio.NewScanner(conn);
	for scanner.Scan() {
		text := scanner.Text()
		// サーバー側に受け取った文字列を表示する
		fmt.Printf("受信メッセージ: %s\n", text)
		conn.Write([]byte("Server: " + text + "\n"))
	}
}