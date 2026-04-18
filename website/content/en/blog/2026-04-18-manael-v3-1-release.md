---
title: "Manael v3.1 is here!"
date: 2026-04-18
description: "Manael v3.1 introduces a post-processing hook for custom handling after image conversion."
---

## Overview

Manael v3.1 introduces a new extension point for applications that embed Manael as a Go library. The new `WithPostProcessor` hook lets you run custom logic after a successful image conversion and before the response is written back to the client.

## New in v3.1: Post-processing Hook

The new hook receives the converted image bytes and can:

- return them unchanged,
- replace them with a different payload, or
- feed them into your own caching or processing pipeline.

This makes it easier to build workflows such as converted-image caching, downstream delivery, or application-specific response handling without forking Manael itself.

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithPostProcessor(func(data []byte) ([]byte, error) {
		return data, nil
	}),
)
```

## Safe fallback behavior

The hook only runs when Manael actually converts an image. If no conversion happens, the hook is skipped. If your hook returns an error, Manael logs the failure and safely falls back to the original upstream response.

## Go 1.24 support has ended

Manael now requires Go 1.25 or later when you build from source or embed it as a library. This change comes with the dependency update to `go.opentelemetry.io/otel/sdk` v1.43.0 for security fixes, which raised the effective minimum supported Go version.

## Learn more

See the new documentation page for the post-processing hook and check the release notes for the full change set in v3.1.
