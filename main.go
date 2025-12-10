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
		args := token[1:]

		switch command {
		case "ADD":
			handleAdd(conn, args)
		case "SEARCH":
			handleSearch(conn, args)
		case "SAVE":
			handleSave(conn)
		default:
			conn.Write([]byte("UNKNOWN COMMAND\n"))
		}
	}
}

func handleAdd(conn net.Conn, args []string) {
	if len(args) < 2 {
		conn.Write([]byte("Error: ADD <id> <vector...>\n"))
		return
	}
	id := args[0]
	vec, err := parseVector(args[1:])
	if err != nil {
		conn.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}
	engine.Add(id, vec)
	conn.Write([]byte("OK\n"))
}

func handleSearch(conn net.Conn, args []string) {
	if len(args) < 1 {
		conn.Write([]byte("Error: SEARCH <vector...>\n"))
		return
	}
	// ベクトルのパース
	query, err := parseVector(args)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}
	// 検索実行（上位３件を表示するように設定してみる）
	results, err := engine.Search(query, 3)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}
	if len(results) == 0 {
		conn.Write([]byte("No results found.\n"))
		return
	}
	// 結果を一行ずつ送信
	var output string
	for _, r := range results {
		output += fmt.Sprintf("ID: %s, Score: %.4f\n", r.ID, r.Score)
	}
	conn.Write([]byte(output))
}

func handleSave(conn net.Conn) {
	err := engine.Save("vectors.json")
	if err != nil {
		conn.Write([]byte("Error: Save failed - " + err.Error() + "\n"))
		return
	}
	conn.Write([]byte("Saved!\n"))
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