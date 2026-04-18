---
title: Image Resizing
weight: 40
---

Manael can resize images on-the-fly using the `w`, `h`, and `fit` query parameters.

## Enabling resizing {#enabling}

Image resizing is disabled by default and must be explicitly opted in by setting the `MANAEL_ENABLE_RESIZE` environment variable to `true`:

```console
MANAEL_ENABLE_RESIZE=true manael -http=:8080 -upstream_url=http://localhost:9000
```

When `MANAEL_ENABLE_RESIZE` is not set to `true`, the `w`, `h`, and `fit` query parameters are silently ignored and the image is converted at its original dimensions.

## Usage {#usage}

When resizing is enabled, append `w` and/or `h` query parameters to the image URL to specify the target dimensions:

```console
curl -sI -H "Accept: image/webp" "http://localhost:8080/image.jpg?w=800&h=600"
```

| Parameter | Description |
| --------- | ----------- |
| `w`       | Target width in pixels (positive integer). |
| `h`       | Target height in pixels (positive integer). |
| `fit`     | Resize mode: `cover`, `contain`, or `scale-down`. Defaults to `contain`. |

## Aspect ratio {#aspect-ratio}

Manael preserves the aspect ratio of the original image when only one dimension is specified.

- **Width only (`?w=300`):** The image is scaled so that its width equals 300 px. The height is adjusted automatically to maintain the original aspect ratio.
- **Height only (`?h=300`):** The image is scaled so that its height equals 300 px. The width is adjusted automatically to maintain the original aspect ratio.
- **Both dimensions (`?w=300&h=300`):** The behavior is controlled by the `fit` parameter:
  - `contain` (default): The image is scaled to fit within the target box while preserving the aspect ratio. The resulting image may be smaller than the requested dimensions.
  - `cover`: The image is scaled and cropped to fill the entire target box while preserving the aspect ratio.
  - `scale-down`: Like `contain`, but the image is never enlarged beyond its original dimensions.

## Security limits {#security-limits}

To prevent excessive resource usage, you can restrict the dimensions that clients may request.

### Maximum dimensions

Set `MANAEL_MAX_RESIZE_WIDTH` and `MANAEL_MAX_RESIZE_HEIGHT` to cap the values accepted for `w` and `h`. Requests that exceed the configured limit are rejected with a `400 Bad Request` response.

```console
MANAEL_ENABLE_RESIZE=true \
MANAEL_MAX_RESIZE_WIDTH=2048 \
MANAEL_MAX_RESIZE_HEIGHT=2048 \
manael -http=:8080 -upstream_url=http://localhost:9000
```

### Allowed dimension lists

For stricter control, `MANAEL_ALLOWED_WIDTHS` and `MANAEL_ALLOWED_HEIGHTS` accept a comma-separated list of permitted pixel values. When set, only the listed values are accepted for the corresponding parameter; any other value results in a `400 Bad Request` response.

```console
MANAEL_ENABLE_RESIZE=true \
MANAEL_ALLOWED_WIDTHS=320,640,1280 \
MANAEL_ALLOWED_HEIGHTS=240,480,720 \
manael -http=:8080 -upstream_url=http://localhost:9000
```

> **Note:** Image resizing is a CPU- and memory-intensive operation. Enable it only when needed and always configure dimension limits in production environments.
