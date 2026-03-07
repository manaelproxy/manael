---
title: デプロイガイド
weight: 30
---

Manael はステートレスな Docker コンテナとして配布されており、最新のコンテナオーケストレーションプラットフォームへの導入が容易です。このガイドでは代表的なデプロイ方法を説明します。

## Docker Compose {#docker-compose}

以下の `docker-compose.yml` は、Manael をローカルの画像サーバーと一緒に動かす例です。

```yaml
services:
  manael:
    image: ghcr.io/manaelproxy/manael:v3
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

次のコマンドで起動します。

```console
docker compose up
```

Manael はポート `8080` でリッスンし、画像リクエストを `origin` サービスに転送します。`MANAEL_ENABLE_AVIF=true` を設定すると、デフォルトの WebP 変換に加えて AVIF 変換も有効になります。

## Google Cloud Run {#cloud-run}

Cloud Run はトラフィックに応じて自動でスケールするため、Manael との相性が良いプラットフォームです。AVIF 変換は CPU 負荷が高いため、メモリ不足エラーを防ぐために少なくとも 1 GiB のメモリを割り当てることを推奨します。

`gcloud` CLI を使って Cloud Run に Manael をデプロイします。

```console
gcloud run deploy manael \
  --image ghcr.io/manaelproxy/manael:v3 \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated \
  --set-env-vars MANAEL_UPSTREAM_URL=https://storage.example.com,MANAEL_ENABLE_AVIF=true \
  --memory 1Gi \
  --cpu 1
```

| フラグ | 説明 |
| --- | --- |
| `--image` | GHCR で公開されている Manael コンテナイメージ。 |
| `--set-env-vars` | `KEY=VALUE` 形式でカンマ区切りの環境変数を指定します。 |
| `--memory` | コンテナインスタンスに割り当てるメモリ量。AVIF 有効時に OOM が発生する場合は `2Gi` に増やしてください。 |
| `--cpu` | コンテナインスタンスに割り当てる vCPU 数。 |
| `--allow-unauthenticated` | 公開アクセスを許可します。アクセスを制限したい場合はこのフラグを削除してください。 |

デプロイ後、Cloud Run がサービス URL を表示します。変換済み画像をキャッシュするために、その URL の前段に CDN や Cloud CDN を配置することを推奨します。

## Kubernetes {#kubernetes}

以下のマニフェストは、Manael を `Deployment` としてデプロイし、クラスター内に公開する `Service` を定義しています。特に AVIF 変換を有効にしている場合は、OOM を防ぐためにリソースの requests と limits を設定することが重要です。

### Deployment と Service {#deployment-and-service}

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
          image: ghcr.io/manaelproxy/manael:v3
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

### リソース設定の目安 {#resource-guidance}

| ワークロード | メモリ request | メモリ limit | 備考 |
| --- | --- | --- | --- |
| WebP のみ（デフォルト） | `128Mi` | `512Mi` | 低〜中程度のトラフィックに適しています。 |
| WebP + AVIF | `256Mi` | `1Gi` | AVIF エンコードは CPU・メモリ使用量がより大きくなります。 |
| 高トラフィック | `512Mi` | `2Gi` | limit をさらに増やすよりも `replicas` を増やして水平スケールさせることを推奨します。 |

Manael はステートレスなので、負荷への対処には `replicas` を増やす水平スケールが最適です。共有ストレージやセッションアフィニティは不要です。
