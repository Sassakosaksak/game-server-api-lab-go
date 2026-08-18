# Game Server API Lab Go

## Router

ルーティングにはchiを採用する。Go標準の`net/http`と互換性があり、標準的なHTTP処理の理解を保ったまま、ルーティングとミドルウェアを段階的に扱えるため。

## 学習ロードマップ

1. Go標準ライブラリで最小APIを作る
2. DockerでAPIをコンテナ化する
3. Docker ComposeでAPI・PostgreSQL・Redisを動かす
4. kindでローカルKubernetesを学ぶ
5. ECS + Fargateへデプロイする
6. EKSも短期検証し、AWS上のKubernetesを体験する

## 学習ノート

仕組みや今回つまずいた点は、[docs/notes](docs/notes/)に分けて記録する。

## Goの単体確認

```powershell
go test ./...
go vet ./...
```

## Dockerで起動

Docker Desktopを起動した状態で、プロジェクトのルートフォルダからImageを作成する。

```powershell
docker build -t game-server-api-go:local .
```

Player APIはPostgreSQLへの接続設定が必要なため、実行・確認には次のDocker Composeを使う。

## Docker Composeで起動

`.env.example`をコピーして`.env`を作成し、`POSTGRES_PASSWORD`へ開発用パスワードを設定する。`.env`はGitへ登録しない。

```powershell
Copy-Item .env.example .env
```

Compose設定を確認する。

```powershell
docker compose config
```

PostgreSQL、Migration、API、Redisを起動する。Migrationは`players`テーブルを作成したあと終了する。Windows側へ公開するのはAPIの`127.0.0.1:8080`だけで、PostgreSQLとRedisはCompose内部ネットワークだけで待ち受ける。

```powershell
docker compose up --build
```

別のPowerShellで状態とAPI応答を確認する。

```powershell
docker compose ps
curl.exe -i http://127.0.0.1:8080/health
docker compose logs migrate
```

Playerを作成する。

```powershell
$body = @{ name = "sumom" } | ConvertTo-Json -Compress
$createdPlayer = Invoke-RestMethod -Method Post -Uri http://127.0.0.1:8080/players -ContentType "application/json" -Body $body
$createdPlayer
```

`$createdPlayer`へ作成されたPlayerが入り、表示すると`id`・`name`・`createdAt`を確認できる。

```powershell
Invoke-RestMethod -Uri "http://127.0.0.1:8080/players/$($createdPlayer.id)"
```

起動したPowerShellで`Ctrl + C`を押したあと、Containerとネットワークを削除する。

```powershell
docker compose down
```

`docker compose down`ではPostgreSQLの名前付きVolumeは残る。`-v`を付けるとVolumeも削除されるため、データを消したい場合だけ使う。
