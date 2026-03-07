---
title: 画像リサイズ
weight: 5
---

Manael は `w`、`h`、`fit` クエリパラメーターを使用して画像をオンザフライでリサイズできます。

## リサイズを有効にする {#enabling}

画像リサイズはデフォルトで無効になっており、`MANAEL_ENABLE_RESIZE` 環境変数を `true` に設定することで明示的に有効にする必要があります。

```console
MANAEL_ENABLE_RESIZE=true manael -http=:8080 -upstream_url=http://localhost:9000
```

`MANAEL_ENABLE_RESIZE` が `true` に設定されていない場合、`w`、`h`、`fit` クエリパラメーターは無視され、画像は元のサイズのまま変換されます。

## 使い方 {#usage}

リサイズが有効な場合、画像 URL に `w` や `h` クエリパラメーターを追加してターゲットの寸法を指定します。

```console
curl -sI -H "Accept: image/webp" "http://localhost:8080/image.jpg?w=800&h=600"
```

| パラメーター | 説明 |
| ----------- | ---- |
| `w`         | ターゲットの幅（ピクセル、正の整数）。 |
| `h`         | ターゲットの高さ（ピクセル、正の整数）。 |
| `fit`       | リサイズモード: `cover`、`contain`、または `scale-down`。デフォルトは `contain`。 |

## アスペクト比 {#aspect-ratio}

一方の寸法のみを指定した場合、Manael は元の画像のアスペクト比を維持します。

- **幅のみ（`?w=300`）:** 幅が 300 px になるよう画像をスケールします。高さはアスペクト比を維持しながら自動的に調整されます。
- **高さのみ（`?h=300`）:** 高さが 300 px になるよう画像をスケールします。幅はアスペクト比を維持しながら自動的に調整されます。
- **両方の寸法（`?w=300&h=300`）:** 動作は `fit` パラメーターによって制御されます。
  - `contain`（デフォルト）: アスペクト比を維持しながら、ターゲットのボックス内に収まるよう画像をスケールします。結果として得られる画像は指定した寸法より小さくなる場合があります。
  - `cover`: アスペクト比を維持しながら、ターゲットのボックス全体を埋めるよう画像をスケールおよびクロップします。
  - `scale-down`: `contain` と同様ですが、元の寸法を超えて拡大されることはありません。

## セキュリティ制限 {#security-limits}

過剰なリソース使用を防ぐため、クライアントがリクエストできる寸法を制限することができます。

### 最大寸法

`MANAEL_MAX_RESIZE_WIDTH` および `MANAEL_MAX_RESIZE_HEIGHT` を設定すると、`w` と `h` に受け付ける値の上限を設定できます。設定した上限を超えるリクエストは `400 Bad Request` レスポンスで拒否されます。

```console
MANAEL_ENABLE_RESIZE=true \
MANAEL_MAX_RESIZE_WIDTH=2048 \
MANAEL_MAX_RESIZE_HEIGHT=2048 \
manael -http=:8080 -upstream_url=http://localhost:9000
```

### 許可する寸法のリスト

さらに厳密な制御として、`MANAEL_ALLOWED_WIDTHS` と `MANAEL_ALLOWED_HEIGHTS` にカンマ区切りで許可するピクセル値のリストを指定できます。設定した場合、対応するパラメーターに対してリストにある値のみが受け付けられ、それ以外の値は `400 Bad Request` レスポンスで拒否されます。

```console
MANAEL_ENABLE_RESIZE=true \
MANAEL_ALLOWED_WIDTHS=320,640,1280 \
MANAEL_ALLOWED_HEIGHTS=240,480,720 \
manael -http=:8080 -upstream_url=http://localhost:9000
```

> **注意:** 画像リサイズは CPU・メモリ使用量の大きい処理です。必要な場合にのみ有効にし、本番環境では必ず寸法の制限を設定してください。
