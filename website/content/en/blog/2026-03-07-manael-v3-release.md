---
title: "Manael v3.0 is here!"
date: 2026-03-07
description: "Introducing Manael v3.0 with libvips integration, image resizing, and dynamic quality control."
---

## Overview

We are excited to announce the release of **Manael v3.0**! This major release brings significant performance improvements and powerful new features for on-the-fly image transformation.

## Powered by libvips

Manael v3.0 now uses [libvips](https://www.libvips.org/) as its image processing engine via the [`bimg`](https://github.com/h2non/bimg) binding. libvips is a fast, memory-efficient image processing library that delivers **significantly better performance and lower memory usage** compared to the previous implementation, making Manael more capable of handling high-traffic workloads.

The official Docker image includes libvips out of the box and requires no additional setup.

## Image Resizing

Manael v3.0 introduces on-the-fly image resizing via new query parameters:

| Parameter | Description             | Example  |
| --------- | ----------------------- | -------- |
| `w`       | Target width in pixels  | `?w=800` |
| `h`       | Target height in pixels | `?h=600` |

You can use `w` and `h` individually or combine them to control how the image is scaled.

### Security-First: Opt-In via `MANAEL_ENABLE_RESIZE`

To protect your deployment from potential abuse and Denial-of-Service attacks, **image resizing is disabled by default**. You must explicitly opt in by setting the `MANAEL_ENABLE_RESIZE` environment variable:

```sh
MANAEL_ENABLE_RESIZE=true
```

When resizing is disabled, `w` and `h` parameters are silently ignored, keeping the default behavior unchanged for existing deployments.

## Dynamic Quality Control

You can now control the compression quality of the output image on a per-request basis using the `q` query parameter:

```text
https://example.com/image.jpg?w=800&q=75
```

The `q` parameter accepts an integer value from `1` to `100`. A higher value produces better image quality at the cost of a larger file size.

A server-level default quality can also be configured with the `MANAEL_DEFAULT_QUALITY` environment variable, giving operators fine-grained control over the quality/size trade-off across all requests.

## Get Started

Update to Manael v3.0 today and take advantage of faster image processing, flexible resizing, and dynamic quality control. For the full list of changes, see the [CHANGELOG](https://github.com/manaelproxy/manael/releases/tag/v3.0.0).
