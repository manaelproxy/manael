---
title: プロダクションガイド
weight: 5
---

このガイドでは、本番環境で Manael を安全かつ効率的に運用するためのベストプラクティスを説明します。

## キャッシュと CDN {#caching-and-cdn}

Manael はステートレスなプロキシであり、変換済み画像を内部にキャッシュしません。本番環境では、不要な再変換を避けてオリジンサーバーの負荷を軽減するために、Manael の前段に CDN またはキャッシュレイヤーを配置してください。

### `Vary: Accept` ヘッダー {#vary-accept-header}

Manael はすべてのレスポンスに `Vary: Accept` ヘッダーを自動的に付与します。このヘッダーは、CDN や中間キャッシュに対して、同一 URL であっても `Accept` ヘッダーの値が異なる場合は別のコピーを保持するよう指示します。これがないと、CDN が WebP 画像をキャッシュし、JPEG しか対応していないクライアントにも WebP を返してしまう可能性があります。

```http
HTTP/1.1 200 OK
Content-Type: image/webp
Vary: Accept
```

1 つの画像 URL に対して、正しく設定された CDN は次の形式ごとに個別のキャッシュエントリを保持します。

- `image/avif` を受け付けるクライアント（AVIF が有効な場合は AVIF で配信）
- `image/webp` を受け付けるクライアント（WebP で配信）
- それ以外のクライアント（元のフォーマットで配信）

### CDN の設定 {#cdn-configuration}

ほとんどの CDN はデフォルトで `Vary` ヘッダーを尊重しますが、明示的な設定が必要なものもあります。利用する CDN について以下を確認してください。

| CDN | 必要な対応 |
| --- | --- |
| Cloudflare | **Polish** を有効にするか、`Vary: Accept` を尊重するキャッシュルールを設定します。 |
| Fastly | デフォルトで `Vary` ヘッダーを尊重するため、追加設定は不要です。 |
| Amazon CloudFront | 許可ヘッダーのリストに `Accept` を含む**キャッシュポリシー**を作成します。 |
| Google Cloud CDN | **キャッシュモード** `CACHE_ALL_STATIC` を有効にし、`Vary` ヘッダーが転送されることを確認します。 |

設定後は、異なる `Accept` ヘッダーで同一 URL にリクエストを送り、それぞれのバリアントが独立してキャッシュされていることを確認してください。

## セキュリティ {#security}

### アップストリームのロック {#upstream-locking}

Manael は起動時に `-upstream_url` フラグ（または `MANAEL_UPSTREAM_URL` 環境変数）でアップストリームのターゲットを固定します。受信したすべてのリクエストは、その単一のオリジンにのみプロキシされます。リクエスト時に別のホストへ転送先を変更することはできません。

この設計により、Manael がオープンプロキシとして悪用されることを防ぎます。攻撃者は、Manael が内部ネットワークやサードパーティのサービスから任意の URL を取得するようなリクエストを作成することができず、SSRF（Server-Side Request Forgery）と呼ばれる脆弱性クラスを排除します。

**推奨事項：**

- `-upstream_url` には、広いネットワークセグメントではなく、特定のバケットやホストなど、できる限り狭いスコープを設定してください。
- ファイアウォールルールやサービスメッシュの egress ポリシーを使用して、Manael のアウトバウンド通信をアップストリームオリジンのみに制限するネットワーク環境で運用してください。

## リソースのチューニング {#resource-tuning}

画像変換、特に AVIF エンコードは CPU とメモリを多く消費します。適切なリソース制限を設定しないと、1 つの Manael インスタンスがノードのリソースを使い果たし、同じホスト上で動作する他のワークロードに影響を与える可能性があります。

### Kubernetes {#kubernetes}

`Deployment` マニフェストには必ずリソースの `requests` と `limits` を設定してください。

```yaml
resources:
  requests:
    cpu: "250m"
    memory: "256Mi"
  limits:
    cpu: "1000m"
    memory: "1Gi"
```

AVIF 変換（`MANAEL_ENABLE_AVIF=true`）を有効にして高負荷なトラフィックを処理する場合は、メモリの limit を `2Gi` に増やしてください。Pod あたりの limit を上げるよりも、`replicas` を増やして水平スケールさせることを推奨します。

### Amazon ECS {#amazon-ecs}

タスク定義に `cpu` と `memory` を設定してください。

```json
{
  "cpu": "1024",
  "memory": "1024"
}
```

`cpu` は CPU ユニット（1024 = 1 vCPU）、`memory` は MiB で指定します。AVIF ワークロードでは `2048` MiB 以上を推奨します。

### 最大画像サイズ {#maximum-image-size}

Manael は、`MANAEL_MAX_IMAGE_SIZE` 環境変数（デフォルト: 20 MiB）で設定されたサイズを超えるアップストリームのレスポンスを拒否します。この値を小さくすると、リクエストごとのピークメモリ使用量が抑えられ、アップストリームから予期せず大きな画像が返ってきた際の影響を軽減できます。

```console
MANAEL_MAX_IMAGE_SIZE=10485760 manael -http=:8080 -upstream_url=https://storage.example.com
```
