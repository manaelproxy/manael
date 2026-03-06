---
title: Deployment Guide
weight: 4
---

Manael is distributed as a stateless Docker container, making it straightforward to deploy to modern container orchestration platforms. This guide covers common deployment scenarios.

## Docker Compose {#docker-compose}

The following `docker-compose.yml` shows how to run Manael alongside a local image server:

```yaml
services:
  manael:
    image: ghcr.io/manaelproxy/manael:v2.1.0
    ports:
      - "8080:8080"
    environment:
      - MANAEL_UPSTREAM_URL=http://origin:9000
      - MANAEL_ENABLE_AVIF=true
    depends_on:
      - origin

  origin:
    image: nginx:alpine
    volumes:
      - ./images:/usr/share/nginx/html:ro
```

Start the stack with:

```console
docker compose up
```

Manael listens on port `8080` and forwards image requests to the `origin` service. Set `MANAEL_ENABLE_AVIF=true` to enable AVIF conversion in addition to the default WebP conversion.

## Google Cloud Run {#cloud-run}

Cloud Run is a good fit for Manael because it scales to zero and handles traffic spikes automatically. Since AVIF conversion is CPU-intensive, allocate at least 1 GiB of memory to avoid out-of-memory errors.

Deploy Manael to Cloud Run with the `gcloud` CLI:

```console
gcloud run deploy manael \
  --image ghcr.io/manaelproxy/manael:v2.1.0 \
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

## Kubernetes {#kubernetes}

The following manifests deploy Manael as a `Deployment` with a `Service` that exposes it inside the cluster. Apply resource requests and limits to prevent OOM kills, especially when AVIF conversion is enabled.

### Deployment and Service {#deployment-and-service}

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: manael
  labels:
    app: manael
spec:
  replicas: 2
  selector:
    matchLabels:
      app: manael
  template:
    metadata:
      labels:
        app: manael
    spec:
      containers:
        - name: manael
          image: ghcr.io/manaelproxy/manael:v2.1.0
          ports:
            - containerPort: 8080
          env:
            - name: MANAEL_UPSTREAM_URL
              value: "https://storage.example.com"
            - name: MANAEL_ENABLE_AVIF
              value: "true"
          resources:
            requests:
              cpu: "250m"
              memory: "256Mi"
            limits:
              cpu: "1000m"
              memory: "1Gi"
---
apiVersion: v1
kind: Service
metadata:
  name: manael
spec:
  selector:
    app: manael
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```

Apply the manifests with:

```console
kubectl apply -f manael.yaml
```

### Resource guidance {#resource-guidance}

| Workload | Memory request | Memory limit | Notes |
| --- | --- | --- | --- |
| WebP only (default) | `128Mi` | `512Mi` | Suitable for low-to-medium traffic. |
| WebP + AVIF | `256Mi` | `1Gi` | AVIF encoding is more CPU- and memory-intensive. |
| High traffic | `512Mi` | `2Gi` | Scale `replicas` horizontally instead of increasing limits further. |

Manael is stateless, so horizontal scaling (increasing `replicas`) is the preferred way to handle additional load. No shared storage or session affinity is required.
