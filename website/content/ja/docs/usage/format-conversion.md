---
title: フォーマット変換
weight: 30
---

このページでは、クライアントの `Accept` ヘッダーや HTTP メソッドに応じて Manael がどのように画像を変換するかを説明します。

## JPEG を WebP に変換する {#converting-jpeg-to-webp}

Manael はクライアントが `Accept` リクエストヘッダーで WebP のサポートを示している場合、JPEG および PNG 画像を自動的に WebP に変換します。プロキシを起動する以外に設定は不要です。

JPEG 画像を WebP として受け取るには、`Accept` ヘッダーに `image/webp` を含めてリクエストします。

```console
curl -sI -X GET -H "Accept: image/webp" http://localhost:8080/image.jpg
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
curl -sI -X GET -H "Accept: image/avif" http://localhost:8080/image.jpg
```

変換が成功すると、`Content-Type: image/avif` が返されます。

```http
HTTP/1.1 200 OK
Content-Type: image/avif
Vary: Accept
...
```

PNG ソース画像への AVIF 変換は、互換性とパフォーマンスを優先するために意図的に無効化されています。PNG 画像はクライアントが対応している場合に引き続き WebP に変換されます。

## リクエストメソッドごとの挙動 {#request-method-behavior}

Manael は `GET` リクエストに対してのみ画像変換を行います。`HEAD` リクエストやその他の HTTP メソッドは変換されず、アップストリームのレスポンスがそのまま返されます。

コマンドラインから変換結果のヘッダーを確認する場合は `curl -sI -X GET` を使用してください。

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
