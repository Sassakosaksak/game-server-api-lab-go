# Redisキャッシュの基本メモ

## Player取得とキャッシュ

PlayerはPostgreSQLを正しい保存先とし、Redisは取得を速くするための一時的なコピーとして使う。

```text
GET /players/{id}
  → Redisから player:<UUID> を取得
  ├─ HIT  : RedisのPlayerを返す
  ├─ MISS : PostgreSQLからPlayerを取得し、Redisへ保存して返す
  └─ 障害 : PostgreSQLからPlayerを取得して返す
```

今回のキーと値は次の形式。

```text
キー  : player:c70b7a6f-350f-40b8-94ff-2f1fb7997ad0
値    : {"id":"...","name":"redistest","createdAt":"..."}
有効期限: 5分
```

RedisはJSON文字列を保存する。`redis-cli GET`では文字列として表示されるため、JSONのダブルクォートがエスケープ付きで見えることがある。

## Redisの障害時はDBへフォールバックする

Redisはキャッシュなので、停止や接続失敗でPlayer取得API全体を失敗させない。Redisからの読み取り・保存に失敗した場合はログを残し、PostgreSQLの結果を返す。

`docker compose stop redis`でRedisサービスを停止すると、Compose内部DNSは`redis`というサービス名を解決できなくなる場合がある。これはAPIからRedisへ接続できない状態を意図的に作る確認方法になる。

## Redis操作の待機時間

Redis障害時に長く待たないよう、Redis操作には200msの上限を設定している。

```text
個別の待機時間
  DialTimeout  : Redisへの接続待ち
  ReadTimeout  : Redisからの応答待ち
  WriteTimeout : Redisへの書き込み待ち
  PoolTimeout  : Redis接続プールから接続を借りる待ち

操作全体の待機時間
  context.WithTimeout
  → 接続・送信・受信を含め、Redis GET / SET 全体で最大200ms
```

`ContextTimeoutEnabled: true`にすると、go-redis内部も`context.WithTimeout`の期限を尊重する。例えば接続に120ms使った場合、残り80ms以内で読み取りまで終える必要がある。

```text
Redisがすぐ「名前を解決できない」と返す
  → 数msでPostgreSQLへフォールバックする

ネットワークが応答しない
  → 最大200msでPostgreSQLへフォールバックする
```

200msはRedisを200ms待つための目標値ではなく、障害時にそれ以上待たせないための上限値。

## 再試行と接続試行

```go
MaxRetries:    -1
DialerRetries: 1
```

`MaxRetries`は、`GET`などRedisコマンド自体の再試行回数。`-1`はgo-redisのコマンド再試行を無効にする指定。

`DialerRetries`は、Redisへ接続できないときの接続試行回数。`0`は既定の5回になるため、今回のように1回にしたい場合は`1`を明示する。

以前は既定値により、コマンドの再試行と接続試行が重なって長く待っていた。現在は接続1回・コマンド再試行なし・操作全体200msの上限としている。
