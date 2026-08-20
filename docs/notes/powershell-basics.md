# PowerShellの基本メモ

## 変数

```powershell
$createdPlayer = Invoke-RestMethod -Method Post ...
$createdPlayer
$createdPlayer.id
```

`$createdPlayer`はPowerShellの変数。`Invoke-RestMethod`が返したJSONをPowerShellのオブジェクトとして、現在開いているPowerShellのメモリへ保存する。

```text
$createdPlayer
  → PowerShellを閉じるまで使える一時変数

PostgreSQL
  → Playerの正式な保存先

Redis
  → Playerの一時的なキャッシュ
```

PowerShell変数はRedisへ保存されない。PowerShellを閉じると変数は消えるが、PostgreSQL・Redisのデータはそれぞれの保存期間やContainer状態に従って残る。

## 処理時間を測る

```powershell
$elapsed = Measure-Command {
    Invoke-RestMethod -Uri "http://127.0.0.1:8080/players/<ID>"
}

$elapsed.TotalMilliseconds
```

`Measure-Command`は、波括弧の中に書いたPowerShell処理の所要時間を測る。`TotalMilliseconds`は結果をミリ秒で取り出す。
