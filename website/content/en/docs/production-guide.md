---
title: Production Guide
weight: 5
---

This guide covers best practices for running Manael safely and efficiently in production environments.

## Caching and CDN {#caching-and-cdn}

Manael is a stateless proxy and does not cache converted images internally. In production, place a CDN or caching layer in front of Manael to avoid redundant conversions and to reduce load on your origin server.

### The `Vary: Accept` header {#vary-accept-header}

Manael automatically includes the `Vary: Accept` header in every response. This header instructs CDNs and intermediate caches to store separate copies of the same URL for each unique `Accept` header value. Without this, a CDN might serve a cached WebP image to a client that only supports JPEG, or vice versa.

```http
HTTP/1.1 200 OK
Content-Type: image/webp
Vary: Accept
```

For a single image URL, a correctly configured CDN will maintain distinct cache entries for:

- Clients that accept `image/avif` (served as AVIF when enabled)
- Clients that accept `image/webp` (served as WebP)
- All other clients (served in the original format)

### CDN configuration {#cdn-configuration}

Most CDNs respect the `Vary` header by default, but some require explicit configuration. Verify the following for your CDN:

| CDN | Action required |
| --- | --- |
| Cloudflare | Enable **Polish** or configure a Cache Rule that respects `Vary: Accept`. |
| Fastly | Respect `Vary` headers by default; no additional configuration needed. |
| Amazon CloudFront | Create a **Cache Policy** that includes `Accept` in the list of allowed headers. |
| Google Cloud CDN | Enable **cache-mode** `CACHE_ALL_STATIC` and ensure `Vary` headers are forwarded. |

Always test your CDN configuration by requesting the same URL with different `Accept` headers and confirming that each variant is cached independently.

## Security {#security}

### Upstream locking {#upstream-locking}

Manael locks the upstream target at startup using the `-upstream_url` flag (or the `MANAEL_UPSTREAM_URL` environment variable). All incoming requests are proxied exclusively to that single upstream origin. It is not possible to redirect Manael to a different host at request time.

This design prevents Manael from being exploited as an open proxy. Attackers cannot craft requests that cause Manael to fetch arbitrary URLs from internal networks or third-party services, which eliminates the class of vulnerabilities known as Server-Side Request Forgery (SSRF).

**Recommendations:**

- Set `-upstream_url` to the narrowest possible scope (e.g., a specific bucket or host) rather than a broad network segment.
- Run Manael in a network environment where outbound traffic is restricted to the upstream origin only, using firewall rules or a service mesh egress policy.

## Resource tuning {#resource-tuning}

Image conversion—especially AVIF encoding—is CPU- and memory-intensive. Without appropriate resource limits, a single Manael instance can exhaust node resources and affect other workloads running on the same host.

### Kubernetes {#kubernetes}

Always set resource `requests` and `limits` in your `Deployment` manifest:

```yaml
resources:
  requests:
    cpu: "250m"
    memory: "256Mi"
  limits:
    cpu: "1000m"
    memory: "1Gi"
```

Increase the memory limit to `2Gi` when AVIF conversion (`MANAEL_ENABLE_AVIF=true`) is enabled under sustained traffic. Prefer horizontal scaling (increasing `replicas`) over increasing per-pod limits.

### Amazon ECS {#amazon-ecs}

Set `cpu` and `memory` in your task definition:

```json
{
  "cpu": "1024",
  "memory": "1024"
}
```

`cpu` is expressed in CPU units (1024 = 1 vCPU) and `memory` in MiB. For AVIF workloads, a value of `2048` MiB or higher is recommended.

### Maximum image size {#maximum-image-size}

Manael rejects upstream responses that exceed the size configured by the `MANAEL_MAX_IMAGE_SIZE` environment variable (default: 20 MiB). Lowering this value reduces peak memory usage per request and limits the impact of unexpectedly large images from the upstream origin.

```console
MANAEL_MAX_IMAGE_SIZE=10485760 manael -http=:8080 -upstream_url=https://storage.example.com
```
