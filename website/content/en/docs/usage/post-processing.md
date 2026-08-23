---
title: Post-processing Hook
weight: 55
---

Manael v3.1 adds Go API hooks for custom processing after image conversion. Use `WithPostProcessor` when you embed Manael as a library and need to inspect, cache, replace, or augment the converted bytes before they are sent to the client. Use `WithRequestPostProcessor` when that processing also needs request-scoped state.

## When the hook runs {#when-it-runs}

The post-processing hook runs only after Manael has successfully converted an image.

- It receives the converted image bytes.
- It can return the same bytes unchanged.
- It can return replacement bytes if your application needs to modify the response.
- It is not called when Manael passes the upstream response through unchanged.

## Basic example {#basic-example}

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithPostProcessor(func(data []byte) ([]byte, error) {
		// Store converted bytes in your cache, object store, or audit pipeline.
		return data, nil
	}),
)
```

## Request-aware processing {#request-aware-processing}

`WithRequestPostProcessor` receives the request together with the converted bytes. Use it to derive a cache key from the request URL or headers, or to apply processing that depends on the caller.

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithRequestPostProcessor(func(r *http.Request, data []byte) ([]byte, error) {
		cacheKey := r.URL.String() + ":" + r.Header.Get("Accept")
		_ = cacheKey // Use the key with your cache or other request-scoped logic.
		return data, nil
	}),
)
```

When both hooks are configured, `WithRequestPostProcessor` takes precedence and `WithPostProcessor` is not called.

## Response headers for the final payload {#response-headers}

`WithResponseHeaderProcessor` receives the inbound request, response headers, and final converted bytes. Use it to add metadata that depends on the final payload, such as an `ETag`, a cache policy, or an audit header. It cannot replace the payload.

Manael calls this hook after it has selected the final bytes and set `Content-Type`, `Content-Length`, and `Content-Disposition` for them.

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithResponseHeaderProcessor(func(r *http.Request, header http.Header, data []byte) error {
		header.Set("ETag", `"`+strconv.FormatInt(int64(len(data)), 10)+`"`)
		header.Set("Cache-Control", "private, max-age=60")
		return nil
	}),
)
```

## Common use cases {#use-cases}

- Store converted images in an external cache.
- Add application-specific response processing after format conversion.
- Integrate with downstream systems that need the final converted payload.
- Set response metadata derived from the final converted payload.

## Error handling {#error-handling}

If the hook returns an error, Manael logs the failure and falls back to the original upstream response instead of sending a partially processed result. This keeps request handling safe even when custom logic fails.

## Important behavior {#behavior-notes}

- The hook is a library API, not a command-line option or environment variable.
- The hook runs after format conversion, so it sees the final converted payload.
- If no conversion happens for the request, the hook is skipped.
- If request-scoped state is needed, use `WithRequestPostProcessor`.
- Use `WithResponseHeaderProcessor` when custom response headers depend on the final payload.
