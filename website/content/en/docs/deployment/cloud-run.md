---
title: Google Cloud Run
weight: 32
---

Cloud Run is a good fit for Manael because it scales to zero and handles traffic spikes automatically. Since AVIF conversion is CPU-intensive, allocate at least 1 GiB of memory to avoid out-of-memory errors.

Deploy Manael to Cloud Run with the `gcloud` CLI:

```console
gcloud run deploy manael \
  --image ghcr.io/manaelproxy/manael:3 \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated \
  --set-env-vars MANAEL_UPSTREAM_URL=https://storage.example.com,MANAEL_ENABLE_AVIF=true \
  --memory 1Gi \
  --cpu 1
```

| Flag | Description |
| --- | --- |
| `--image` | The Manael container image from GHCR. |
| `--set-env-vars` | Comma-separated `KEY=VALUE` pairs for environment variables. |
| `--memory` | Memory allocated per container instance. Increase to `2Gi` if you experience OOM kills with AVIF. |
| `--cpu` | Number of vCPUs per container instance. |
| `--allow-unauthenticated` | Allows public access. Remove this flag if you want to restrict access. |

After deployment, Cloud Run prints the service URL. Place a CDN or Cloud CDN in front of that URL to cache the converted images.
