# Goテストの基本メモ

## ポートを使わないHTTP APIテスト

`httptest`を使うと、APIサーバーを`ListenAndServe`で起動せずにRouterとHandlerをテストできる。

```text
テストコード
  → 偽のHTTPリクエスト
  → chi Router
  → healthHandler
  → 偽のHTTPレスポンス受け口
  → ステータス・ヘッダー・JSONを確認
```

この方式ではPCのポートを待ち受けず、Windows FirewallやDockerも使わない。

```go
request := httptest.NewRequest(http.MethodGet, "/health", nil)
recorder := httptest.NewRecorder()

newRouter().ServeHTTP(recorder, request)
response := recorder.Result()
```

`httptest.NewRequest`は偽のHTTPリクエストを作る。`httptest.NewRecorder`は、Handlerが`ResponseWriter`へ書いたステータス・ヘッダー・本文をメモリ上に記録するテスト用の受け口。

`newRouter()`は、本番起動とテストで同じルーティング設定を使うために作る。本番側だけルートを変え、テスト側の登録を直し忘れることを防ぐ。

## testing.T

```go
func TestHealth(t *testing.T)
```

`go test`は`Test`から始まる関数を見つけ、テストの進行役である`*testing.T`を渡して呼び出す。`t.Fatalf`は失敗内容を記録し、その時点でそのテストを終了する。

## テスト対象の指定

```powershell
go test -v ./...
go test -v ./cmd/api
go test -v -run '^TestHealth$' ./cmd/api
```

```text
./...      : 現在のフォルダ以下にある全Goパッケージ
./cmd/api  : cmd/apiパッケージだけ
-run       : 実行するテスト名を正規表現で絞る
```

`^TestHealth$`は`TestHealth`という名前に完全一致する。正規表現が一致せずテストが0件でも、`go test`は失敗ではなく`PASS`になるため、`testing: warning: no tests to run`が出ていないか確認する。

```powershell
go test .\cmd\api\main_test.go
```

これはテストファイルだけをコンパイルする指定。`main.go`の`newRouter`や`healthResponse`を含まないため、通常は使わない。テストではファイル名ではなくパッケージのフォルダを指定する。
