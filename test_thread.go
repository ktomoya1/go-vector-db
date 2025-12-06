package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. バイト（Goroutine）を雇って走らせる
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("　バイト： 皿洗い中…", i)
			time.Sleep(500 * time.Millisecond) // 0.5秒休む
		}
	}()

	// 2. 店長（main）も同時に走る
	for i := 0; i < 5; i++ {
		fmt.Println("店長：　レジ打ち中…", i)
		time.Sleep(500 * time.Millisecond) // 0.5秒休む
	}
}