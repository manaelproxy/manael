---
title: Docker Compose
weight: 31
---

Use Docker Compose when you want to run Manael alongside a local image server or a small self-managed deployment.

The following `docker-compose.yml` shows how to run Manael alongside a local image server:

```yaml
services:
  manael:
    image: ghcr.io/manaelproxy/manael:3
    ports:
      - "8080:8080"
    environment:
      - MANAEL_UPSTREAM_URL=http://origin:9000
      - MANAEL_ENABLE_AVIF=true
    depends_on:
      - origin

  origin:
    image: nginx:alpine
    volumes:
      - ./images:/usr/share/nginx/html:ro
```

Start the stack with:

```console
docker compose up
```

Manael listens on port `8080` and forwards image requests to the `origin` service. Set `MANAEL_ENABLE_AVIF=true` to enable AVIF conversion in addition to the default WebP conversion.