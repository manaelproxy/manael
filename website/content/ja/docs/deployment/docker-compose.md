---
title: Docker Compose
weight: 31
---

Docker Compose は、Manael をローカルの画像サーバーと一緒に動かす場合や、小規模な単一ホスト構成で手早く試したい場合に向いています。

以下の `docker-compose.yml` は、Manael をローカルの画像サーバーと一緒に動かす例です。

```yaml
services:
  manael:
    image: ghcr.io/manaelproxy/manael:3
    ports:
      - "8080:8080"
    environment:
      - MANAEL_UPSTREAM_URL=http://origin
      - MANAEL_ENABLE_AVIF=true
    depends_on:
      - origin

  origin:
    image: nginx:alpine
    volumes:
      - ./images:/usr/share/nginx/html:ro
```

次のコマンドで起動します。

```console
docker compose up
```

Manael はポート `8080` でリッスンし、画像リクエストを `origin` サービスに転送します。`MANAEL_ENABLE_AVIF=true` を設定すると、デフォルトの WebP 変換に加えて AVIF 変換も有効になります。
