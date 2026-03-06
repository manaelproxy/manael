---
title: "Manael v3.0 Release"
date: 2026-03-06
draft: true
---

## Overview

We are pleased to announce the release of **Manael v3.0**. This is a major release that introduces significant changes to image processing, dependency management, and query parameter support.

## Migration to libvips

Manael v3.0 migrates the underlying image processing library to [libvips](https://www.libvips.org/) via the [`bimg`](https://github.com/h2non/bimg) binding. libvips is a fast, memory-efficient image processing library that provides better performance and broader format support compared to the previous implementation.

### Dynamic Linking Requirement

As a result of this migration, **libvips must be installed on the host system** at runtime. Manael is no longer a fully statically linked binary.

**On Debian/Ubuntu:**

```sh
sudo apt-get install -y libvips-dev
```

**On Alpine Linux:**

```sh
apk add vips-dev
```

The official Docker image includes libvips and requires no additional setup.

## New Query Parameters

Manael v3.0 introduces new query parameters for on-the-fly image resizing:

| Parameter | Description                              | Example        |
|-----------|------------------------------------------|----------------|
| `w`       | Target width in pixels                   | `?w=800`       |
| `h`       | Target height in pixels                  | `?h=600`       |
| `fit`     | Resize fit mode (`cover`, `contain`, etc.) | `?fit=cover` |

These parameters can be combined to control how the image is scaled and cropped.

## Breaking Changes

- The binary now requires `libvips` to be present on the host at runtime.
- Minimum Go version has been updated. Please check the `go.mod` file for the current requirement.

## How to Upgrade

1. Install `libvips` on your host or update your container image.
2. Replace the Manael binary with the v3.0 release.
3. Review the query parameter documentation if you use URL-based image transformations.

For the full list of changes, see the [CHANGELOG](https://github.com/manaelproxy/manael/blob/main/CHANGELOG.md).
