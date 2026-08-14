# Game Server API Lab Go

## 学習ロードマップ

1. Go標準ライブラリで最小APIを作る
2. DockerでAPIをコンテナ化する
3. Docker ComposeでAPI・PostgreSQL・Redisを動かす
4. kindでローカルKubernetesを学ぶ
5. ECS + Fargateへデプロイする
6. EKSも短期検証し、AWS上のKubernetesを体験する

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
