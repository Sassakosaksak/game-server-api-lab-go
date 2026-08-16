# VS CodeでGoを書くときのメモ

## Go拡張とgopls

VS Codeでは、公式のGo拡張（提供元: Go Team at Google、ID: `golang.go`）と`gopls`がGoコードの解析を担当する。

`gopls`が起動していると、補完、引数情報、ホバー表示、定義ジャンプ、参照検索、エラー表示が使える。右下にGoのバージョンと`⚡`が表示されているかを確認する。

```text
Ctrl + Shift + P
  → Go: Locate Configured Go Tools
```

このコマンドで`gopls`の有無と、プロジェクトの`go.mod`を見ているかを確認できる。

`gopls`が無い場合は、次から導入する。

```text
Ctrl + Shift + P
  → Go: Install/Update Tools
  → goplsを選ぶ
```

## 主な操作

```text
Ctrl + Space         : 候補を表示する
Ctrl + Shift + Space : 関数呼び出し中に引数情報を表示する
F12 または Ctrl + クリック : 定義へ移動する
Alt + F12            : 定義をその場で表示する
Shift + F12          : 呼び出し元・参照箇所を探す
Alt + ←              : 定義ジャンプ前の場所へ戻る
Alt + →              : 戻る操作の後に進む
```

同じフォルダかつ同じ`package`のGoファイルは、importなしで関数や型を共有できる。別フォルダは通常別パッケージなので、importして`パッケージ名.関数名`の形で使う。
