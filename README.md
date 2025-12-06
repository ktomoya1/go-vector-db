# Go Vector DB Project

Go言語で実装する、学習用のインメモリ・ベクトルデータベース。
TCPソケット通信を使い、Redisライクなコマンドで操作する。

## プロジェクトの目的
- Go言語の並行処理（Goroutine, Channel, Mutex）の理解
- ネットワークプログラミング（TCP/IP, Socket）の理解
- データ構造とアルゴリズム（Map, Vector Search）の実装

## フェーズ分け計画

### Phase 1: 基本的なKVS（現在ここ！）
- [x] TCPサーバーの立ち上げ (Listen, Accept)
- [x] 複数クライアントの同時接続 (Goroutine)
- [x] クライアントとのメッセージ交換 (Echo)
- [ ] データの保存機能 (In-Memory Map)
- [ ] 排他制御 (Sync.Mutex)
- [ ] コマンド実装 (SET, GET, DEL)

### Phase 2: ベクトル検索エンジンの実装
- [ ] 文字列だけでなく、ベクトルデータ([0.1, 0.5...])を保存可能にする
- [ ] コサイン類似度計算の実装
- [ ] 類似データの検索コマンド (SEARCH)

### Phase 3: 永続化と最適化 (Optional)
- [ ] データをファイルに保存 (Save/Load)
- [ ] インデックス化による高速化

## 技術スタック
- 言語: Go
- ライブラリ: 標準ライブラリのみ (net, bufio, sync, fmt)