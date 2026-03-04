---
title: Usage
weight: 3
---

## Starting the proxy {#starting-the-proxy}

Start Manael by specifying the address to listen on and the upstream server URL:

```console
$ manael -http=:8080 -upstream_url=http://localhost:9000
```

## Converting JPEG to WebP {#converting-jpeg-to-webp}

Manael converts JPEG and PNG images to WebP automatically when the client signals support for WebP via the `Accept` request header. No configuration is needed beyond starting the proxy.

To request a JPEG image and receive it as WebP, include `image/webp` in the `Accept` header:

```console
$ curl -sI -H "Accept: image/webp" http://localhost:8080/image.jpg
```

When conversion succeeds, the response includes `Content-Type: image/webp`:

```
HTTP/1.1 200 OK
Content-Type: image/webp
Vary: Accept
...
```

If the client does not include `image/webp` in the `Accept` header, the original image is returned unchanged.
