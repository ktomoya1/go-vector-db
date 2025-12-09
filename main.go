package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
	scanner := bufio.NewScanner(conn)
	
	for scanner.Scan() {
		line := scanner.Text()
		token := strings.Fields(line)

		if len(token) == 0 {
			continue
		}

		command := token[0]

		switch command {
		case "ADD":
			if len(token) < 3 {
				conn.Write([]byte("Error: ADD <id> <vector...>\n"))
				continue
			}
			id := token[1]
			vec, err := parseVector(token[2:])
			if err != nil {
				conn.Write([]byte("Error: " + err.Error() + "\n"))
				continue
			}
			engine.Add(id, vec)
			conn.Write([]byte("OK\n"))
		default:
			conn.Write([]byte("UNKNOWN COMMAND\n"))
		}
	}
}

func parseVector(args []string) (Vector, error) {
	var v Vector
	for _, s := range args {
		// 文字列を64bit浮動小数点数に変換
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", s)
		}
		v = append(v, val)
	}
	return v, nil
}