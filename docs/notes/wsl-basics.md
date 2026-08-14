# WSLの基本メモ

## WSL 2とUbuntu

```text
Windows
  └─ WSL 2
      └─ Linuxカーネル
          ├─ Ubuntu
          └─ docker-desktop
```

WSLはWindows上でLinuxを動かす仕組み。WSL 2は軽量な仮想環境の中でLinuxカーネルを動かす方式で、DockerのLinuxコンテナにも使われる。

UbuntuはWSL 2の上で動かすLinuxディストリビューション。`docker-desktop`はDocker Desktopが内部用に作る環境であり、普段のLinux学習には使わない。

Docker Desktopだけを使うなら、Ubuntuの導入は必須ではない。今回はLinuxのファイル操作、パッケージ管理、AWS上のLinuxサーバーの理解も学ぶため、自分用のUbuntuを導入した。

## パス

Linuxのフォルダは`/`を頂点に持つ。

```text
/home/sumom
  → Ubuntu内の自分のホームフォルダ

/mnt/c/Users/sumom/GameServerApi-GO
  → Windowsの C:\Users\sumom\GameServerApi-GO
```

`/mnt/c`配下はWindows側のCドライブをUbuntuから見ている場所。今のプロジェクトはWindows側に置いたままDockerを学ぶ。

Linux側で大量のファイル操作をする開発では`/home/sumom`配下の方が速い場合があるが、今はプロジェクトを移動しない。

Ubuntu内のファイルは、Windowsのエクスプローラーからも次のパスで見られる。

```text
\\wsl.localhost\Ubuntu\home\sumom
```

Ubuntuで`explorer.exe .`を実行すると、Ubuntuで現在いるフォルダをWindowsのエクスプローラーで開ける。
