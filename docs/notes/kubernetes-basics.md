# Kubernetesの基本メモ

## kindとYAMLの`kind`は別物

`kind`はKubernetes IN Dockerの略。Docker ContainerをKubernetes Nodeとして使い、ローカルクラスタを作るツール。

```text
Windows PC
  └─ Docker Desktop
      └─ kind Node Container
          └─ Kubernetes Pod
```

YAMLの`kind: Deployment`などは、作るKubernetesリソースの種類を表す項目であり、kindツールとは別物。

## Cluster、Node、Namespace、Pod、Container、Service

NodeとNamespaceは別軸。

```text
Cluster
├─ Node A、Node B
│   → Podを実際に動かすコンピュータや仮想マシン
│
└─ Namespace: game-server
    → API、DB、Redisなどを論理的に整理するグループ
```

PodはNamespaceへ所属し、SchedulerによってNodeへ配置される。

```text
PostgreSQL Pod: postgres-0
  → Namespace: game-server
  → Node: game-server-local-control-plane
  → Container: postgres
```

ServiceはPodやContainerの親ではない。同じNamespace内で、ラベルが一致するPodへの通信入口を作る。

```text
Service selector: app=redis
  → Pod label: app=redis
  → 一致したRedis Podへ通信を送る
```

## Workload Controller

DeploymentとStatefulSetは、どちらもPodを維持するWorkload Controller。

```text
Deployment
  → 交換可能なPodを必要数維持する
  → APIやRedisなど、Pod固有の名前・保存領域を意識しない処理に使う

StatefulSet
  → 順序・安定したPod名・保存領域を意識してPodを維持する
  → PostgreSQLなど、状態を持つDBに使う
```

今回のPostgreSQL Pod名は`postgres-0`。`replicas: 1`なので1個だけ作られる。

Nodeを増やしてもPodは自動で増えない。Pod数はDeploymentやStatefulSetの`replicas`で決まる。

```text
kindの role: worker
  → Worker Nodeを作る

spec.replicas: 2
  → Podを2個作る
```

複数Nodeで配置先を決めたい場合は、Nodeへ自分でラベルを付け、Pod側の`nodeSelector`で一致するNodeを指定する。

```text
Node A: workload=api
Node B: workload=db

PostgreSQL Pod
  nodeSelector: workload=db
  → Node Bへだけ配置できる
```

`workload`は予約語ではなく、自分で決めたラベル名。

## Serviceとポート

APIがPostgreSQLへ接続する場合、接続先はService名とServiceポート。

```text
API Pod
  → postgres-internal:5432
  → PostgreSQL PodのIP:5432
  → PostgreSQL Container
```

```yaml
# Service
metadata:
  name: postgres-internal
spec:
  ports:
    - port: 5432
      targetPort: postgres-tcp

# PostgreSQL Container
ports:
  - name: postgres-tcp
    containerPort: 5432
```

`port`はAPIなどの接続元がServiceへ指定する番号。`targetPort`は転送先のContainerポートで、番号またはContainerポート名を指定できる。`postgres-tcp`はTCPを有効にする設定ではなく、`targetPort`から参照するための名前札。

Containerがポートを1つだけ使う場合は、番号を直接書いてもよい。

```yaml
port: 5432
targetPort: 5432
```

Containerが複数ポートを使う場合は、名前を付けるとServiceの接続先が分かりやすい。

```yaml
# API Container
ports:
  - name: http
    containerPort: 8080
  - name: metrics
    containerPort: 9090

# 通常API用Service
targetPort: http

# 監視用Service
targetPort: metrics
```

PostgreSQL用の`postgres-internal`はHeadless Service。

```yaml
clusterIP: None
```

Headless Serviceは専用のService IPを作らず、DNSが対象PodのIPを返す。今回のようなStatefulSetでは、`postgres-0`のような安定したPod名と組み合わせる。

Headless Serviceは、複数DBが同じ5432番を使えるようにするための設定ではない。接続先はIPアドレスとポート番号の組み合わせで決まるため、PodのIPが別なら同じ5432番を使える。

```text
PostgreSQL Pod A: 10.244.0.7:5432
PostgreSQL Pod B: 10.244.0.9:5432
  → IPが違うため、同じ5432番でも区別できる
```

Headless Serviceを使うと、StatefulSetが作る個別PodをDNS名で指定できる。

```text
postgres-internal:5432
  → 今回はpostgres-0だけなので、そのPodへ接続する

postgres-0.postgres-internal:5432
postgres-1.postgres-internal:5432
  → 複数PostgreSQL Podがある場合に、個別Podを指定できる
```

通常はAPIから個別Pod名へ接続せず、`postgres-internal:5432`を使う。個別Pod名はPostgreSQLレプリケーションなど、DB同士が相手を区別して通信するときに使う。

Redis用の`redis-internal`は通常Service。Service専用のCluster IPを持ち、対象Podへ通信を振り分ける。

Pod内のContainerは通常、同じPod IPを共有する。`containerPort`はポートを開放する命令ではなく、Containerが使うポートを名前付きで宣言する項目。実際に待ち受けるのはPostgreSQLやRedisのプロセス。

## StorageClass、PVC、Volume

```text
StorageClass: standard
  → 保存領域を作る方式

PVC: postgres-data
  → PostgreSQL用に保存領域が欲しい、という要求

PV
  → StorageClassがPVCに対して用意する実際の保存領域

Pod volumes
  → PVCをPodで使う

Container volumeMounts
  → Podの保存領域をContainer内のフォルダへ見せる
```

```text
PVC postgres-data
  → Podの volumes: postgres-data
  → Containerの volumeMounts
  → /var/lib/postgresql/data
```

`mountPath`はNode上のフォルダではなく、PostgreSQL Container内のフォルダ。PostgreSQLはここへDBデータを書き込む。

`ReadWriteOnce`は、同時に1 Nodeから読み書きするアクセスモード。PostgreSQL Podが1個の今回に合う。

```text
Podだけ再作成
  → 同じPVCを再びマウントできる
  → DBデータは残る

kindのNode自体が消える
  → local-pathの保存領域も失われ得る
```

kindの`standard`はNode側のローカル保存領域を使う。AWS EBSなどの外部ディスクのように、別Nodeへディスクを付け替える検証は今回の1 Node kind環境では行わない。

複数PostgreSQL Podを作る場合、同じDBデータフォルダを共有して書き込まない。通常はPodごとに別PVCを持たせ、DBレプリケーションで同期する。

## Secretと`.env`

Docker Composeはローカルの`.env`を読んでContainerの環境変数へ渡せる。

Kubernetesクラスタは自分のPCの`.env`を直接読めない。パスワードはSecretとしてクラスタへ登録し、必要なPodだけが環境変数やファイルとして参照する。

```text
ローカル .env
  → kubectl create secret
  → Kubernetes Secret
  → PostgreSQL Podの POSTGRES_PASSWORD
```

Secretの実値はGitへ登録しない。

## Probe

`livenessProbe`は、応答できないPodをKubernetesが再作成するための確認。

`readinessProbe`は、Serviceから新しい通信を渡してよいかの確認。失敗してもPodは再起動せず、ReadyになるまでServiceの接続先から外す。

Podが`Running`でも、アプリケーションが接続を受け付けられるとは限らない。

```text
PostgreSQL Containerのプロセスは起動済み
  → PodはRunning
  → ただし起動処理・復旧処理中で接続をまだ受け付けない
  → pg_isreadyは失敗
  → READYは0/1
  → Serviceから新しい通信を流さない

起動処理が完了
  → pg_isreadyが成功
  → READYは1/1
  → Serviceが通信を流せる
```

APIの場合は、APIプロセス自体は起動していても、DBの接続先が誤っている・DBが停止中・必要なMigrationが未適用、といった理由で実際の処理を受け付けられないことがある。そうした依存先も確認するreadinessProbeを作ることがある。

```text
initialDelaySeconds
  → Pod起動後、最初の確認まで待つ秒数

periodSeconds
  → 以後の確認間隔

timeoutSeconds
  → 1回の確認コマンドを待つ上限。未指定時は既定1秒
```

## `kubectl`の確認コマンド

```powershell
kubectl get pods,services,pvc -n game-server -o wide
```

`-n game-server`は、`game-server` Namespace内を対象にする指定。`-o wide`はPod IPやNode名なども表示する。

```powershell
kubectl exec -n game-server postgres-0 -c postgres -- pg_isready -U gameuser -d gamedb
```

```text
exec              : Container内でコマンドを実行する
-n game-server    : 対象Namespace
postgres-0        : 対象Pod
-c postgres       : 対象Container
--                : 以後はContainer内で実行するコマンド
```

Serviceは通信の入口でありContainerを持たないため、`kubectl exec`の対象にはできない。
