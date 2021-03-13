---
title: インストールガイド
---

## Docker を使う {#using-the-docker}

Manael は [Docker](https://www.docker.com/) で動かすことを推奨しています。Manael の Docker イメージは [Docker Hub](https://hub.docker.com/) で公開されています。

Manael を Docker を使って動かす場合は `docker pull manael/manael:latest` コマンドで取得してください。Docker を使うことによって既存の環境に不必要なファイルを増やさずに最新版の Manael が使えるようになります。

## ビルド済みバイナリを使う {#using-a-built-binary}

64 ビット版の GNU/Linux を対象にしてビルドされた Manael をダウンロードできます。

### 1. ディレクトリを作る {#1-create-a-working-directory}

まず Manael をインストールする際にダウンロードしたファイルを展開するためのディレクトリを作ります。

```console
$ mkdir manael
$ cd manael
```

### 2. ダウンロード {#2-download}

[リリースページ](https://github.com/manaelproxy/manael/releases)から最新版の Manael (`manael_1.x.y_Linux_x86_64.tar.gz`) をダウンロードして 1. で作ったディレクトリに展開します。

```console
$ wget https://github.com/manaelproxy/manael/releases/download/v1.x.y/manael_1.x.y_Linux_x86_64.tar.gz
$ tar xf manael_1.x.y_Linux_x86_64.tar.gz
```

### 3. インストール {#3-install}

ファイルをコピーするために `install` コマンドを利用します。`cp` コマンドや `mv` コマンドでも同様の作業はできますが、`install` コマンドを使うことによって適切な実行権限が実行ファイルに与えられます。

```console
$ sudo install manael /usr/local/bin
```

## ソースコードからビルド {#build-from-a-source-code}

Manael のソースコードは [GitHub](https://github.com/manaelproxy/manael) にホストされています。Manael は [Go](https://golang.org/) で書かれていて、`go` コマンドを使って簡単にビルドできます。

```console
$ go build -o manael cmd/manael/main.go
```
