# Goの基本メモ

## Goモジュール

`go.mod`があるフォルダから下が、1つのGoモジュールになる。

```go
module github.com/Sassakosaksak/game-server-api-lab-go
```

`module`は、自分のコードをimportするときの基準となる名前。Gitの接続先は`.git/config`の`origin`であり、`go.mod`ではない。

## 実行とビルド

```text
go run ./cmd/api
  → 一時的にビルドして起動する
  → プロジェクト内には実行ファイルを残さない

go build -o bin/api.exe ./cmd/api
  → ./cmd/api を材料にしてビルドする
  → bin/api.exe へ実行ファイルを出力する
```

`-o`の後ろは出力先、最後の`./cmd/api`はビルド対象のGoパッケージ。

## HTTP APIの最小構成

```go
package main
```

`package main`と`func main()`の組み合わせは、起動できるアプリケーションを表す。

```go
type healthResponse struct {
	Status string `json:"status"`
}
```

`struct`は複数の値をまとめる型。`Status`はGoコード内のフィールド名、``json:"status"``はJSONへ変換するときのキー名。

Goでは先頭が大文字の名前は公開される。`encoding/json`は公開フィールドをJSONとして出力するため、`Status`を大文字にしている。

## エラー確認

```go
if err := json.NewEncoder(w).Encode(value); err != nil {
	log.Printf("書き込みに失敗しました: %v", err)
}
```

`:=`は「変数の宣言と代入」を同時に行う短縮記法。`;`の左で`err`を作り、右で`err != nil`を判定している。

`nil`は「値がない」状態。`error`、ポインタ、スライス、mapなどで使える。`int`や`string`は`nil`にできない。

## ポインタ

```go
count := 10
countPointer := &count
*countPointer = 20
```

`&count`は`count`が置かれている場所を取得する。`*countPointer`は、その場所にある値を読む、または書き換える。

```text
T   : 値そのものの型
*T  : Tへのポインタ型
&x  : xへのポインタを作る
```

`server := &http.Server{}`のようにポインタで扱うと、同じサーバー設定を複数の変数で共有できる。
