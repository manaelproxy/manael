---
title: Kubernetes
weight: 33
---

Use Kubernetes when you need explicit control over replicas, resource limits, and service exposure inside a cluster.

The following manifests deploy Manael as a `Deployment` with a `Service` that exposes it inside the cluster. Apply resource requests and limits to prevent OOM kills, especially when AVIF conversion is enabled.

## Deployment and Service {#deployment-and-service}

Save the following manifest to a file named `manael.yaml`:

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
                  image: ghcr.io/manaelproxy/manael:3
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

Apply the manifest with:

```console
kubectl apply -f manael.yaml
```

## Resource guidance {#resource-guidance}

| Workload            | Memory request | Memory limit | Notes                                                               |
| ------------------- | -------------- | ------------ | ------------------------------------------------------------------- |
| WebP only (default) | `128Mi`        | `512Mi`      | Suitable for low-to-medium traffic.                                 |
| WebP + AVIF         | `256Mi`        | `1Gi`        | AVIF encoding is more CPU- and memory-intensive.                    |
| High traffic        | `512Mi`        | `2Gi`        | Scale `replicas` horizontally instead of increasing limits further. |

Manael is stateless, so horizontal scaling (increasing `replicas`) is the preferred way to handle additional load. No shared storage or session affinity is required.
