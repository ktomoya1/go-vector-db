# Simple Vector Search Engine in Go 🚀

Go言語の標準ライブラリのみで実装した、インメモリ・ベクトル検索エンジンのプロトタイプです。
TCPソケット通信を通じて、Redisのようなコマンドラインインターフェースで操作できます。

## 概要
RAG（Retrieval-Augmented Generation）の中核技術である「ベクトル検索」の仕組みを深く理解するために、**コサイン類似度計算をゼロからGo言語で実装**しました。
ブラックボックスになりがちな「AIが似ている文章を探す仕組み（高次元ベクトルの内積計算）」を、自分の手で再現・検証するためのプロジェクトです。

## 🏗 アーキテクチャ

関心の分離（Separation of Concerns）を意識し、通信層と計算ロジック層を明確に分けて設計しました。

```mermaid
classDiagram
    class Main_Controller {
        <<main.go>>
        +TCP接続ハンドリング
        +コマンド解析 (Parser)
    }
    class VectorEngine {
        <<engine.go>>
        +データの保存 (Map)
        +排他制御 (Mutex)
        +コサイン類似度計算
    }
    Main_Controller --> VectorEngine : 利用