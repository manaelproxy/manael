---
title: Kubernetes
weight: 33
---

Kubernetes は、レプリカ数、リソース制限、クラスター内での公開方法を明示的に制御したい場合に向いています。

以下のマニフェストは、Manael を `Deployment` としてデプロイし、クラスター内に公開する `Service` を定義しています。特に AVIF 変換を有効にしている場合は、OOM を防ぐためにリソースの requests と limits を設定することが重要です。

## Deployment と Service {#deployment-and-service}

以下のマニフェストを `manael.yaml` という名前のファイルに保存します。

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

次のコマンドでマニフェストを適用します。

```console
kubectl apply -f manael.yaml
```

## リソース設定の目安 {#resource-guidance}

| ワークロード            | メモリ request | メモリ limit | 備考                                                                                 |
| ----------------------- | -------------- | ------------ | ------------------------------------------------------------------------------------ |
| WebP のみ（デフォルト） | `128Mi`        | `512Mi`      | 低〜中程度のトラフィックに適しています。                                             |
| WebP + AVIF             | `256Mi`        | `1Gi`        | AVIF エンコードは CPU・メモリ使用量がより大きくなります。                            |
| 高トラフィック          | `512Mi`        | `2Gi`        | limit をさらに増やすよりも `replicas` を増やして水平スケールさせることを推奨します。 |

Manael はステートレスなので、負荷への対処には `replicas` を増やす水平スケールが最適です。共有ストレージやセッションアフィニティは不要です。
