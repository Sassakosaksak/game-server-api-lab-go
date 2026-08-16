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
COPY go.mod go.sum ./
COPY cmd ./cmd
```

`COPY`は左がコピー元、右がコピー先。`WORKDIR /src`の後の`./`はImage内の`/src`を指す。

`go.mod`は、このAPIがどのモジュールのどの版を使うかを宣言するファイル。例えば`github.com/go-chi/chi/v5 v5.3.1`は「chiのv5.3.1を使う」という意味。

`go.sum`は、取得した依存パッケージの内容を検証するためのチェックサム記録。`chi`では次の2行が作られる。

```text
github.com/go-chi/chi/v5 v5.3.1 h1:...
github.com/go-chi/chi/v5 v5.3.1/go.mod h1:...
```

1行目はchi本体のソースコード全体、2行目はchi自身の`go.mod`の内容に対するチェックサム。外部パッケージを使うGo APIでは、Docker内でも同じ依存関係を再現・検証できるよう両方をコピーする。

```text
go.sumの要約
  1行目 : 取り込んだモジュール本体のコードが改ざんされていないか確認する
  2行目 : 取り込んだモジュール自身のgo.modが改ざんされていないか確認する
```

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
docker run --rm --name game-server-api-go -p 127.0.0.1:8080:8080 game-server-api-go:local
```

`-p 127.0.0.1:8080:8080`は、Windowsの自PCだけが使える8080番ポートをコンテナ側の8080番へ転送する。

```text
curl.exe http://localhost:8080/health
  → Windowsの8080番
  → コンテナ内のGo APIの8080番
```

ポート公開はWindows Firewallの確認対象になり得る。

## ビルド元とタグ

```powershell
docker build -t game-server-api-go:local .
```

最後の`.`はビルド材料にするフォルダ。プロジェクトのルートで実行すると、そのフォルダ内の標準名`Dockerfile`を使う。

ImageタグはDocker Desktop全体で共有される名前札であり、プロジェクトごとに分かれてはいない。タグ名が同じなら、最後にそのタグを付けてビルドしたImageを指す。ビルド失敗時には、前に成功したImageがそのタグのまま残り得る。

## Docker Desktopの開始と停止

```powershell
docker desktop status
docker desktop start
docker desktop stop
```

`docker desktop stop`はDocker Engineを停止するため、起動中のLinux Containerも停止する。`--rm`付きで起動したContainerは終了時に自動削除される。Imageは残るので、Docker Desktopを再起動すれば再び`docker run`できる。

Docker Desktopのアンインストールは終了とは別で、ローカルのContainer・Image・VolumeなどのDockerデータを削除する操作。

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
