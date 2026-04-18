---
title: Post-processing Hook
weight: 55
---

Manael v3.1 adds a Go API hook for custom processing after image conversion. Use `WithPostProcessor` when you embed Manael as a library and need to inspect, cache, replace, or augment the converted bytes before they are sent to the client.

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

## Common use cases {#use-cases}

- Store converted images in an external cache.
- Add application-specific response processing after format conversion.
- Integrate with downstream systems that need the final converted payload.

## Error handling {#error-handling}

If the hook returns an error, Manael logs the failure and falls back to the original upstream response instead of sending a partially processed result. This keeps request handling safe even when custom logic fails.

## Important behavior {#behavior-notes}

- The hook is a library API, not a command-line option or environment variable.
- The hook runs after format conversion, so it sees the final converted payload.
- If no conversion happens for the request, the hook is skipped.
