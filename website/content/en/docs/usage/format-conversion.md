---
title: Format Conversion
weight: 30
---

This page explains how Manael converts images based on the client `Accept` header and request method.

## Converting JPEG to WebP {#converting-jpeg-to-webp}

Manael converts JPEG and PNG images to WebP automatically when the client signals support for WebP via the `Accept` request header. No configuration is needed beyond starting the proxy.

To request a JPEG image and receive it as WebP, include `image/webp` in the `Accept` header:

```console
curl -sI -X GET -H "Accept: image/webp" http://localhost:8080/image.jpg
```

When conversion succeeds, the response includes `Content-Type: image/webp`:

```http
HTTP/1.1 200 OK
Content-Type: image/webp
Vary: Accept
...
```

If the client does not include `image/webp` in the `Accept` header, the original image is returned unchanged.

## Converting images to AVIF {#converting-to-avif}

Manael can also convert JPEG images to AVIF. AVIF conversion is disabled by default and must be enabled with the `MANAEL_ENABLE_AVIF` environment variable:

```console
MANAEL_ENABLE_AVIF=true manael -http=:8080 -upstream_url=http://localhost:9000
```

When AVIF is enabled and the client includes `image/avif` in the `Accept` header, Manael converts eligible images to AVIF:

```console
curl -sI -X GET -H "Accept: image/avif" http://localhost:8080/image.jpg
```

A successful conversion returns `Content-Type: image/avif`:

```http
HTTP/1.1 200 OK
Content-Type: image/avif
Vary: Accept
...
```

AVIF conversion is intentionally disabled for PNGs to prioritize compatibility and performance. PNGs are still converted to WebP when the client supports it.

## Request method behavior {#request-method-behavior}

Manael only converts images for `GET` requests. `HEAD` requests and other HTTP methods are passed through to the upstream server unchanged.

When testing from the command line, use `curl -sI -X GET` to fetch the headers of a converted image.

## Content-Disposition header updates {#content-disposition-updates}

When Manael converts an image to a different format, it automatically updates the filename extension in the `Content-Disposition` response header to match the new format. For example, if an upstream response includes:

```http
Content-Disposition: attachment; filename="photo.jpg"
```

and the image is converted to WebP, Manael updates the header to:

```http
Content-Disposition: attachment; filename="photo.webp"
```

The same update applies when converting to AVIF (`.avif` extension).
