---
title: Quality Control
weight: 6
---

Manael allows clients to control the encoding quality of converted images via the `q` query parameter.

## Usage {#usage}

Append the `q` query parameter to the image URL to specify the desired encoding quality. The value must be an integer in the range `1` to `100`, where higher values produce better quality at larger file sizes:

```console
curl -sI -H "Accept: image/webp" "http://localhost:8080/image.jpg?q=80"
```

## Valid range {#valid-range}

The `q` parameter accepts integer values between `1` and `100` (inclusive). Values outside this range are clamped automatically:

- Values below `1` are treated as `1`.
- Values above `100` are treated as `100`.

## Format-specific quality {#format-specific}

You can target a specific output format using the `format:value` syntax. Multiple format-specific values can be combined in a single `q` parameter using a comma-separated list:

```console
# Set quality only for WebP
curl -sI -H "Accept: image/webp" "http://localhost:8080/image.jpg?q=webp:75"

# Set quality for WebP and AVIF independently
curl -sI -H "Accept: image/avif,image/webp" "http://localhost:8080/image.jpg?q=webp:80,avif:50"

# Set a universal quality with a WebP-specific override
curl -sI -H "Accept: image/webp" "http://localhost:8080/image.jpg?q=70,webp:85"
```

| Syntax | Description |
| ------ | ----------- |
| `?q=80` | Universal quality applied to any output format. |
| `?q=webp:75` | Quality applied only when the image is converted to WebP. |
| `?q=avif:50` | Quality applied only when the image is converted to AVIF. |
| `?q=webp:80,avif:50` | Format-specific values combined in one parameter. |

## Default behavior {#default-behavior}

When the `q` parameter is omitted, Manael uses built-in per-format quality defaults:

| Format | Default quality |
| ------ | --------------- |
| WebP   | 90              |
| AVIF   | 60              |

### Server-level default quality

Operators can override the built-in defaults by setting the `MANAEL_DEFAULT_QUALITY` environment variable. This value applies to all formats unless the client provides a format-specific override via the `q` parameter:

```console
MANAEL_DEFAULT_QUALITY=75 manael -http=:8080 -upstream_url=http://localhost:9000
```

The resolution order when determining the encoding quality for a given request is:

1. Format-specific quality from the `q` query parameter (e.g., `?q=webp:80`).
2. Universal quality from the `q` query parameter (e.g., `?q=80`).
3. Server-level default from `MANAEL_DEFAULT_QUALITY`.
4. Built-in per-format default (90 for WebP, 60 for AVIF).
