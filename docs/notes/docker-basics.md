# Dockerの基本メモ

## Dockerfile、Image、Container

```text
Dockerfile + Goソースコード
  → docker build
  → Image
  → docker run
  → Container
```

- Dockerfile: Imageの作り方を書くレシピ
- Image: 実行に必要なものを梱包した完成品
- Container: Imageを実際に起動した実行中の環境

## Image名とタグ

```text
game-server-api-go:local
├─ game-server-api-go : Image名
└─ local              : タグ
```

タグはImageの版や用途を区別する名前札。

```text
game-server-api-go:local
game-server-api-go:v1
game-server-api-go:v2
```

同じImage名でもタグを指定して別バージョンを起動できる。同じタグで再ビルドすると、タグは新しいImageを指す。

## Dockerfileの2段階ビルド

```text
golang:1.26.5
  → GoコンパイラでLinux用の /out/api を作る

alpine:3.22
  → /out/api だけをコピーして実行する
```

Goの実行用ImageにはGoランタイムを入れない。C#でいうSDK ImageでpublishしてからASP.NET Runtime Imageで動かす構成に近い。

```dockerfile
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
```

`COPY`は左がコピー元、右がコピー先。`WORKDIR /src`の後の`./`はImage内の`/src`を指す。

```dockerfile
RUN ... go build -o /out/api ./cmd/api
```

`-o /out/api`は完成した実行ファイルの出力先。`./cmd/api`はビルド対象のGoコードの場所。

```dockerfile
EXPOSE 8080
```

これはコンテナ内で8080番を使う宣言。外部公開はしない。実際にWindows側とつなぐのは`docker run -p`。

## コンテナのポート公開

```powershell
docker run --rm --name game-server-api-go -p 8080:8080 game-server-api-go:local
```

`-p 8080:8080`は、Windows側の8080番ポートをコンテナ側の8080番へ転送する。

```text
curl.exe http://localhost:8080/health
  → Windowsの8080番
  → コンテナ内のGo APIの8080番
```

ポート公開はWindows Firewallの確認対象になり得る。

## Docker Composeとの関係

```yaml
api:
  build: .
```

`build: .`は、Composeファイルのあるフォルダを材料に、標準名の`Dockerfile`を使ってAPI Imageを作る指定。

```yaml
db:
  image: postgres:16
```

`image:`は既存Imageを取得してコンテナを起動する指定。ComposeはAPI、DB、Redisなど複数コンテナをまとめて定義・起動する。
