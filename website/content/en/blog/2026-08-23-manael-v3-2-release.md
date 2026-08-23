---
title: "Manael v3.2 is here!"
date: 2026-08-23
description: "Manael v3.2 adds request-aware post-processing and response-header customization based on the final payload."
---

## Overview

Manael v3.2 expands the post-processing hooks available to applications that embed Manael as a Go library. You can now inspect the incoming request while processing converted images and customize response headers based on the final converted payload.

## Request-aware post-processing

The new `WithRequestPostProcessor` receives the incoming `*http.Request` alongside the converted image bytes. It is useful when you need to generate cache keys from the request URL or `Accept` header, or apply processing specific to the caller.

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithRequestPostProcessor(func(r *http.Request, data []byte) ([]byte, error) {
		cacheKey := r.URL.String() + ":" + r.Header.Get("Accept")
		_ = cacheKey // Use this for request-scoped work such as caching.
		return data, nil
	}),
)
```

When both `WithRequestPostProcessor` and the existing `WithPostProcessor` are configured, the request-aware hook takes precedence.

## Response headers based on the final payload

`WithResponseHeaderProcessor` receives the incoming request, response headers, and final converted bytes. Use it to set metadata derived from the final payload, such as `ETag`, cache policy, or audit headers. This hook cannot replace the payload itself.

Manael invokes the hook after it has determined the final bytes and set `Content-Type`, `Content-Length`, and `Content-Disposition`.

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithResponseHeaderProcessor(func(r *http.Request, header http.Header, data []byte) error {
		header.Set("ETag", `"`+strconv.FormatInt(int64(len(data)), 10)+`"`)
		header.Set("Cache-Control", "private, max-age=60")
		return nil
	}),
)
```

## Safe fallback behavior

These hooks only run when Manael actually converts an image. If no conversion happens, they are skipped. If either hook returns an error, Manael logs the failure and safely falls back to the original upstream response instead of returning a partially processed result.

## Learn more

See the post-processing hook documentation and check the release notes for the full change set in v3.2.
