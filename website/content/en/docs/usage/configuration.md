---
title: Configuration
weight: 60
---

In addition to command-line flags, Manael supports configuration via environment variables.

| Environment variable   | Equivalent flag | Description                                                                                 |
| ---------------------- | --------------- | ------------------------------------------------------------------------------------------- |
| `MANAEL_UPSTREAM_URL`  | `-upstream_url` | URL of the upstream image server.                                                           |
| `PORT`                 | `-http`         | Port number to listen on (used as `:<PORT>`).                                               |
| `MANAEL_ENABLE_AVIF`   | —               | Set to `true` to enable AVIF conversion.                                                    |
| `MANAEL_ENABLE_RESIZE` | —               | Set to `true` to enable on-the-fly image resizing via `w`, `h`, and `fit` query parameters. |

The `-upstream_url` flag takes precedence over `MANAEL_UPSTREAM_URL` when both are provided.
