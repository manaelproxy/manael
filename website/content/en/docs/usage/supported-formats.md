---
title: Supported Formats
weight: 10
---

This page describes the image formats that Manael can accept as input and the output formats it produces.

## Conversion table {#conversion-table}

| Input format   | Output format                        | Notes                                              |
| -------------- | ------------------------------------ | -------------------------------------------------- |
| JPEG           | WebP / AVIF (if enabled)             |                                                    |
| PNG            | WebP                                 | AVIF output is intentionally disabled for PNG      |
| APNG           | Pass-through (original APNG)         | Conversion is skipped to preserve animation data   |
| Animated GIF   | Pass-through (original GIF)          | Conversion is skipped to preserve animation data   |
| Static GIF     | WebP                                 | Added in v3.0.0                                    |

## Notes {#notes}

### AVIF for JPEG only {#avif-jpeg-only}

AVIF conversion is supported for JPEG source images only. PNG images are intentionally excluded from AVIF conversion to preserve transparency and maintain broad compatibility. PNG images are still converted to WebP when the client supports it.

### Pass-through for animated formats {#pass-through-animated}

Manael automatically detects animated images — both Animated PNG (APNG) and Animated GIF — and returns them to the client unchanged. This behavior prevents losing animation data that would occur if an animated image were converted to a single static frame in WebP or AVIF. The original file is returned regardless of the `Accept` header sent by the client.
