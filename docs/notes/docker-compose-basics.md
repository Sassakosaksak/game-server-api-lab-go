# Docker Composeの基本メモ

## Composeでまとめて起動する

`compose.yaml`は、複数のContainerをまとめて定義するファイル。

```text
docker compose up --build
  → api、db、redisをまとめて起動する
```

今回の構成では、Windows側に公開するのはAPIだけ。

```text
Windowsのcurl
  → 127.0.0.1:8080
  → api Container

api Container
  → db:5432    : PostgreSQL
  → redis:6379 : Redis
```

`db`と`redis`はComposeのサービス名。Composeが作る内部ネットワークでは、IPアドレスを固定せずサービス名で接続先を指定できる。

## depends_onとhealthcheck

```yaml
depends_on:
  db:
    condition: service_healthy
```

`depends_on`は起動順序の依存関係。`condition: service_healthy`を付けると、PostgreSQLのヘルスチェックが成功してからAPIを起動する。

```text
pg_isready     : PostgreSQLが接続受付可能か確認する
redis-cli ping : RedisがPONGを返せるか確認する
```

## .envとパスワード

`.env`はComposeが変数を読み込むローカル設定ファイル。`POSTGRES_PASSWORD`のような秘密値は`.env`へ置き、Gitへ登録しない。

`.env.example`はGitへ登録する雛形であり、実際のパスワードは書かない。

## PostgreSQLのVolume

```yaml
volumes:
  - postgres-data:/var/lib/postgresql/data
```

`postgres-data`はDockerが管理する名前付きVolume。PostgreSQLのデータをContainerの外へ保存する。

```text
docker compose down
  → Containerとネットワークを削除
  → 名前付きVolumeは残る

docker compose down -v
  → 名前付きVolumeも削除
  → PostgreSQLのデータも消える
```
