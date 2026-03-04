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
