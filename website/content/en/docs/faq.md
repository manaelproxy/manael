---
title: FAQ
weight: 90
---

## Why doesn't Manael convert PNGs to AVIF? {#png-avif}

AVIF encoding for PNGs is intentionally disabled to prioritize compatibility, encoding speed, and alpha channel handling. PNGs will still be converted to WebP when the client supports it.

See the [Supported Formats]({{< relref "supported-formats" >}}) page for more details on the conversion table.

## Can I use an Amazon S3 bucket URL as the `-upstream_url`? {#s3-upstream}

Yes, as long as the S3 bucket is publicly accessible or accessible from the network where Manael is running.

## Does Manael cache the converted images? {#caching}

No. Manael acts purely as a stateless processing proxy. Please use a CDN in front of Manael for caching.
