# Simple Vector Search Engine in Go

## 概要
RAG（Retrieval-Augmented Generation）の中核技術である「ベクトル検索」の仕組みを深く理解するために、**コサイン類似度計算をゼロからGo言語で実装した**検索エンジンのプロトタイプです。

## 開発の動機
「AIがどうやって似ている文章を探しているのか？」という疑問を持ち、Azure OpenAI Serviceなどの裏側にある数学的なロジック（高次元ベクトルの内積計算）を、ブラックボックスのままにせず自分の手で再現したくて開発しました。

## 技術的なこだわり
1. **アルゴリズムの自作**: 外部ライブラリに頼らず、コサイン類似度（Cosine Similarity）の計算式を実装し、数学的な理解を深めました。
2. **並行処理を意識した設計**: 将来的に大量の検索リクエストを捌くことを想定し、Go言語標準の `sync.RWMutex` を使用してスレッドセーフな設計にしました。
3. **計算量**: 線形探索 $O(N)$ ですが、データの正規化やメモリアロケーションを意識しています。

## 今後の展望
* OpenAI API等と連携し、実際のテキストデータをEmbeddingして検索できるようにする。
* データ量が増えた際に、Goroutineを使った並列計算で検索速度を維持する。

classDiagram
    %% 定義
    class Main_Controller {
        <<main.go>>
        +main()
        +handleConnection(net.Conn)
        -parseVector(args) Vector
    }

    class VectorEngine {
        <<engine.go>>
        -Store map[string]Vector
        -Mu sync.RWMutex
        +NewVectorEngine()
        +Add(key, vec)
        +Search(query, limit)
        -CosineSimilarity(v1, v2)
    }

    class Vector {
        <<Type Alias>>
        []float64
    }

    %% 関係性
    Main_Controller --> VectorEngine : 利用 (Uses)
    VectorEngine *-- Vector : 保持 (Contains)
    Main_Controller ..> Vector : 生成 (Creates)

    note for Main_Controller "・TCP接続の受付\n・コマンド解析 (ADD/SEARCH)\n・結果の整形・応答"
    note for VectorEngine "・データの保存 (Map)\n・排他制御 (Mutex)\n・数学的計算 (Cosine)"
