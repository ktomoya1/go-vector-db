package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

type Vector []float64

// コサイン類似度を計算する関数
func CosineSimilarity(v1, v2 Vector) (float64, error) {
	// 配列の長さ（次元）が違ったら角度を計算できないためエラーにする
	if len(v1) != len(v2) {
		return 0, fmt.Errorf("vector lengths do not match")
	}

	// 内積
	var dotProduct float64
	// ベクトルAの長さ
	var normA float64
	// ベクトルBの長さ
	var normB float64

	// 次元（配列の長さ）の数だけ繰り返す
	for i := 0; i < len(v1); i++ {
		dotProduct += v1[i] * v2[i]
		normA += v1[i] * v1[i]
		normB += v2[i] * v2[i]
	}

	// 0除算を防ぐ
	if normA == 0 || normB == 0 {
		return 0, fmt.Errorf("zero vector found")
	}

	// cosθ = (A ・ B) / (||A|| * ||B||)
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}

// 検索結果の表示用
type SearchResult struct {
	ID    string
	Score float64
}

// 検索エンジンクラス
type VectorEngine struct {
	// 辞書の定義。キーはstring型。値はVector型。
	data map[string]Vector
	mu   sync.RWMutex
}

// エンジンクラスのインスタンス化
func NewVectorEngine() *VectorEngine {
	return &VectorEngine{
		data: make(map[string]Vector),
	}
}

func (ve *VectorEngine) Add(id string, v Vector) {
	// データに書き込んでいる間はロックする
	ve.mu.Lock()
	defer ve.mu.Unlock()
	ve.data[id] = v
}

// queryには検索文字のベクトル数値、limitには表示件数を入れる
func (ve *VectorEngine) Search(query Vector, limit int) ([]SearchResult, error) {
	// 検索している間は検索データをロックする
	ve.mu.RLock()
	defer ve.mu.RUnlock()

	// 裏でメモリを取るが、appendの性質上、SearchResultのサイズは0にしておく
	results := make([]SearchResult, 0, len(ve.data))

	// idには検索文字列、vecにはその文字列のベクトルが入る
	// ve.dataのサイズだけ繰り返す
	for id, vec := range ve.data {
		score, err := CosineSimilarity(query, vec)
		if err != nil {
			continue
		}
		// 現在の検索データが検索対象とどれだけ近い意味を持っているかの数値を追加する
		results = append(results, SearchResult{ID: id, Score: score})
	}

	// スコア列を大きい順に並べる
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// 決めた表示件数よりスコア列が大きかったら、表示件数分切り取る
	if len(results) > limit {
		return results[:limit], nil
	}
	return results, nil
}