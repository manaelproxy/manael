---
title: 使い方
weight: 3
---

## プロキシを起動する {#starting-the-proxy}

待ち受けるアドレスとアップストリームサーバーの URL を指定して Manael を起動します。

```console
manael -http=:8080 -upstream_url=http://localhost:9000
```

## JPEG を WebP に変換する {#converting-jpeg-to-webp}

Manael はクライアントが `Accept` リクエストヘッダーで WebP のサポートを示している場合、JPEG および PNG 画像を自動的に WebP に変換します。プロキシを起動する以外に設定は不要です。

JPEG 画像を WebP として受け取るには、`Accept` ヘッダーに `image/webp` を含めてリクエストします。

```console
curl -sI -H "Accept: image/webp" http://localhost:8080/image.jpg
```

変換が成功すると、レスポンスに `Content-Type: image/webp` が含まれます。

```http
HTTP/1.1 200 OK
Content-Type: image/webp
Vary: Accept
...
```

クライアントが `Accept` ヘッダーに `image/webp` を含めない場合、元の画像がそのまま返されます。

## 画像を AVIF に変換する {#converting-to-avif}

Manael は JPEG 画像を AVIF に変換することもできます。AVIF 変換はデフォルトで無効になっており、`MANAEL_ENABLE_AVIF` 環境変数で有効にする必要があります。

```console
MANAEL_ENABLE_AVIF=true manael -http=:8080 -upstream_url=http://localhost:9000
```

AVIF が有効で、クライアントが `Accept` ヘッダーに `image/avif` を含めている場合、Manael は対象画像を AVIF に変換します。

```console
curl -sI -H "Accept: image/avif" http://localhost:8080/image.jpg
```

変換が成功すると、`Content-Type: image/avif` が返されます。

```http
HTTP/1.1 200 OK
Content-Type: image/avif
Vary: Accept
...
```

> **注意:** PNG ソース画像への AVIF 変換は、互換性とパフォーマンスを優先するために意図的に無効化されています。PNG 画像はクライアントが対応している場合に引き続き WebP に変換されます。

## アニメーション PNG (APNG) のパススルー {#apng-pass-through}

Manael はアニメーション PNG (APNG) ファイルを自動的に検出し、変換せずにそのままクライアントへ送信します。APNG を静止画の WebP や AVIF に変換するとアニメーションデータが失われるため、`Accept` ヘッダーの内容に関わらず元の APNG がそのまま返されます。

## Content-Disposition ヘッダーの更新 {#content-disposition-updates}

Manael が画像を別フォーマットに変換した際、レスポンスの `Content-Disposition` ヘッダーに含まれるファイル名の拡張子を新しいフォーマットに合わせて自動的に更新します。たとえば、アップストリームのレスポンスに次のヘッダーが含まれる場合：

```http
Content-Disposition: attachment; filename="photo.jpg"
```

WebP に変換されると、Manael はヘッダーを次のように更新します：

```http
Content-Disposition: attachment; filename="photo.webp"
```

AVIF に変換した場合も同様に拡張子が `.avif` に更新されます。

## 高度な設定 {#advanced-configuration}

### 環境変数 {#environment-variables}

Manael はコマンドラインフラグに加えて、環境変数による設定もサポートしています。

| 環境変数                | 対応するフラグ    | 説明                                                          |
| ---------------------- | --------------- | ------------------------------------------------------------- |
| `MANAEL_UPSTREAM_URL`  | `-upstream_url` | アップストリーム画像サーバーの URL。                            |
| `PORT`                 | `-http`         | 待ち受けるポート番号（`:<PORT>` の形式で使用されます）。         |
| `MANAEL_ENABLE_AVIF`   | —               | `true` に設定すると AVIF 変換が有効になります。                 |

`-upstream_url` フラグと `MANAEL_UPSTREAM_URL` の両方が指定されている場合、フラグが優先されます。
