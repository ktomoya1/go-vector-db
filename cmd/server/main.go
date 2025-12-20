package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"vector-db/pkg/vector"
)

var engine *vector.VectorEngine

func main() {
	engine = vector.NewVectorEngine()

	loadErr := engine.Load("vectors.json")
	if loadErr != nil {
		fmt.Printf("Warning: Failed to load data (starting with empty DB): %v\n", loadErr)
	}

	fmt.Println("Vector Engine initialized!")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server started. Waiting for connections on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		fmt.Println("New connection established!")
		
		conn.Write([]byte("Hello! You are connected to My-Vector-DB!\n"))
		
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// メモリ上に一度まとめて読み込んでから処理を行うbufioが高速
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
		// case "BENCHMARK":
		// 	handleBenchmark(conn, args)
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

	query, err := parseVector(args)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}

	// Default limit set to 3
	results, err := engine.Search(query, 3)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}
	if len(results) == 0 {
		conn.Write([]byte("No results found.\n"))
		return
	}

	var output string
	for _, r := range results {
		output += fmt.Sprintf("ID: %s, Score: %.4f\n", r.ID, r.Score)
	}
	conn.Write([]byte(output))
}

func parseVector(args []string) (vector.Vector, error) {
	var v vector.Vector
	for _, s := range args {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", s)
		}
		v = append(v, val)
	}
	return v, nil
}

func handleSave(conn net.Conn) {
	err := engine.Save("vectors.json")
	if err != nil {
		conn.Write([]byte("Error: Save failed - " + err.Error() + "\n"))
		return
	}
	conn.Write([]byte("Saved!\n"))
}

// func handleBenchmark(conn net.Conn, args []string) {
// }
