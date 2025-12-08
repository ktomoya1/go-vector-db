package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var engine *VectorEngine

func main() {
	engine = NewVectorEngine()

	engine.Add("test_doc", Vector{0.1, 0.2, 0.3})
	fmt.Println("ベクトルエンジン初期化完了!")

	// 1. サーバーを立ち上げる（ポート8080番で待ち受け）
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("サーバー起動失敗:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("サーバーが起動しました！クライアントからの接続を待っています...")

	for {
		// 2. 誰かが接続してくるのを待つ（ここでプログラムは一時停止します）
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接続エラー:", err)
			continue
		}

		// 3. 接続が来たら、処理を開始する
		fmt.Println("新しい接続がありました！")
		
		// 接続相手に挨拶を送る
		conn.Write([]byte("Hello! You are connected to My-Vector-DB!\n"))
		
		// 相手からのメッセージを読み取る準備
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	
	// 入力を読み取るスキャナーを作る
	scanner := bufio.NewScanner(conn)
	
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("受信メッセージ: %s\n", text)
		
		// オウム返し（Echo）をする
		conn.Write([]byte("Server received: " + text + "\n"))
	}
}