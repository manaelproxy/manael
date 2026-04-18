---
title: Google Cloud Run
weight: 32
---

Cloud Run はトラフィックに応じて自動でスケールするため、Manael との相性が良いプラットフォームです。AVIF 変換は CPU 負荷が高いため、メモリ不足エラーを防ぐために少なくとも 1 GiB のメモリを割り当てることを推奨します。

`gcloud` CLI を使って Cloud Run に Manael をデプロイします。

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

| フラグ                    | 説明                                                                                                    |
| ------------------------- | ------------------------------------------------------------------------------------------------------- |
| `--image`                 | GHCR で公開されている Manael コンテナイメージ。                                                         |
| `--set-env-vars`          | `KEY=VALUE` 形式でカンマ区切りの環境変数を指定します。                                                  |
| `--memory`                | コンテナインスタンスに割り当てるメモリ量。AVIF 有効時に OOM が発生する場合は `2Gi` に増やしてください。 |
| `--cpu`                   | コンテナインスタンスに割り当てる vCPU 数。                                                              |
| `--allow-unauthenticated` | 公開アクセスを許可します。アクセスを制限したい場合はこのフラグを削除してください。                      |

デプロイ後、Cloud Run がサービス URL を表示します。変換済み画像をキャッシュするために、その URL の前段に CDN や Cloud CDN を配置することを推奨します。
