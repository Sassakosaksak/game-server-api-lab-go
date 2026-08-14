# Game Server API Lab Go

## 学習ロードマップ

1. Go標準ライブラリで最小APIを作る
2. DockerでAPIをコンテナ化する
3. Docker ComposeでAPI・PostgreSQL・Redisを動かす
4. kindでローカルKubernetesを学ぶ
5. ECS + Fargateへデプロイする
6. EKSも短期検証し、AWS上のKubernetesを体験する

## 学習ノート

仕組みや今回つまずいた点は、[docs/notes](docs/notes/)に分けて記録する。

## ローカル起動

```powershell
go run ./cmd/api
```

別のPowerShellで、次のコマンドを実行してAPIを確認する。

```powershell
curl http://localhost:8080/health
```

期待する応答：

```json
{"status":"ok"}
```

## Dockerで起動

Docker Desktopを起動した状態で、プロジェクトのルートフォルダからImageを作成する。

```powershell
docker build -t game-server-api-go:local .
```

作成したImageからコンテナを起動する。

```powershell
docker run --rm --name game-server-api-go -p 8080:8080 game-server-api-go:local
```

別のPowerShellで、APIとコンテナログを確認する。

```powershell
curl.exe http://localhost:8080/health
docker logs game-server-api-go
```

コンテナを起動したPowerShellで`Ctrl + C`を押すと、コンテナは停止して自動削除される。
