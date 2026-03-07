---
title: Usage
weight: 20
---

## Starting the proxy {#starting-the-proxy}

Start Manael by specifying the address to listen on and the upstream server URL:

```console
manael -http=:8080 -upstream_url=http://localhost:9000
```

## Converting JPEG to WebP {#converting-jpeg-to-webp}

Manael converts JPEG and PNG images to WebP automatically when the client signals support for WebP via the `Accept` request header. No configuration is needed beyond starting the proxy.

To request a JPEG image and receive it as WebP, include `image/webp` in the `Accept` header:

```console
curl -sI -H "Accept: image/webp" http://localhost:8080/image.jpg
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
curl -sI -H "Accept: image/avif" http://localhost:8080/image.jpg
```

A successful conversion returns `Content-Type: image/avif`:

```http
HTTP/1.1 200 OK
Content-Type: image/avif
Vary: Accept
...
```

> **Note:** AVIF conversion is intentionally disabled for PNG source images to prioritize compatibility and performance. PNG images are still converted to WebP when the client supports it.

## On-the-fly image resizing {#image-resizing}

Manael can resize images on-the-fly using the `w`, `h`, and `fit` query parameters. Resizing is disabled by default and must be enabled with the `MANAEL_ENABLE_RESIZE` environment variable:

```console
MANAEL_ENABLE_RESIZE=true manael -http=:8080 -upstream_url=http://localhost:9000
```

When resizing is enabled, clients can request a specific size by appending query parameters to the image URL:

```console
curl -sI -H "Accept: image/webp" "http://localhost:8080/image.jpg?w=800&h=600&fit=cover"
```

| Parameter | Description |
| --------- | ----------- |
| `w`       | Target width in pixels (positive integer). |
| `h`       | Target height in pixels (positive integer). |
| `fit`     | Resize mode: `cover`, `contain`, or `scale-down`. Defaults to `contain`. |

When `MANAEL_ENABLE_RESIZE` is not set to `true`, resize parameters (`w`, `h`, `fit`) are silently ignored and the image is converted without resizing.

> **Note:** Image resizing is a CPU- and memory-intensive operation. Enable it only when needed and consider setting `MANAEL_MAX_RESIZE_WIDTH` and `MANAEL_MAX_RESIZE_HEIGHT` to limit resource usage.

## Animated PNG (APNG) pass-through {#apng-pass-through}

Manael automatically detects Animated PNG (APNG) files and passes them through to the client without conversion. This prevents losing animation data that would occur if an APNG were converted to a static WebP or AVIF frame. The original APNG is returned unchanged regardless of the `Accept` header.

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

## Advanced configuration {#advanced-configuration}

### Environment variables {#environment-variables}

In addition to command-line flags, Manael supports configuration via environment variables.

| Environment variable   | Equivalent flag  | Description                                          |
| ---------------------- | ---------------- | ---------------------------------------------------- |
| `MANAEL_UPSTREAM_URL`  | `-upstream_url`  | URL of the upstream image server.                    |
| `PORT`                 | `-http`          | Port number to listen on (used as `:<PORT>`).        |
| `MANAEL_ENABLE_AVIF`   | —                | Set to `true` to enable AVIF conversion.             |
| `MANAEL_ENABLE_RESIZE` | —                | Set to `true` to enable on-the-fly image resizing via `w`, `h`, and `fit` query parameters. |

The `-upstream_url` flag takes precedence over `MANAEL_UPSTREAM_URL` when both are provided.
