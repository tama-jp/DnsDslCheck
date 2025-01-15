# DnsDslCheck

- 本番配布用のバイナリにコンパイル (linuxはクロスコンパイルが難しいのでやっていない)

```shell
wails build -platform darwin/universal
wails build -platform windows
```

- 開発の準備が整っているかを診断

```shell
wails doctor
```


- ライブ開発"モードで実行

```shell
wails dev -v 2 
```

- Wails CLIのバージョンをアップデート

```shell
wails update
```

- 現在のCLIバージョンを出力するだけのコマンド

```shell
wails version
```
