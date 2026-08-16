# Goの基本メモ

## Goモジュール

`go.mod`があるフォルダから下が、1つのGoモジュールになる。

```go
module github.com/Sassakosaksak/game-server-api-lab-go
```

`module`は、自分のコードをimportするときの基準となる名前。Gitの接続先は`.git/config`の`origin`であり、`go.mod`ではない。

## 依存パッケージの追加

Goで外部パッケージを使う流れは次の通り。

```text
1. Goコードへimportを書く
2. go get または go mod tidyを実行する
3. Goがgo.modとgo.sumを更新する
4. go.modとgo.sumをGitへコミットする
```

```powershell
go get github.com/go-chi/chi/v5
go mod tidy
```

`go get`は、C#の`dotnet add package`に近い。`go.mod`へ使うモジュールと版を追加し、モジュールを取得する。

`go mod tidy`は、コードのimportを見て依存関係を整理する。importされているのに`go.mod`に無いモジュールは追加し、どこからも使われていないモジュールは削除する。必要な`go.sum`も追加・整理する。

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

## 書式指定

`log.Printf`や`fmt.Printf`、`t.Fatalf`では、文章内の`%`から始まる書式指定へ後ろの値を埋め込める。

```go
t.Fatalf("ステータスコード = %d, want %d", actualStatus, expectedStatus)
```

```text
%d : 整数を10進数で表示する
%q : 文字列をダブルクォート付きで表示する
%s : 文字列をそのまま表示する
%v : 型に応じた基本表示をする。errorの内容表示などに使う
```

`%q`は空文字や余計な空白も見分けやすいため、文字列を比較するテストで使いやすい。
