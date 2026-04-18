---
title: 使い始める
weight: 10
---

待ち受けるアドレスとアップストリームサーバーの URL を指定して Manael を起動します。

```console
manael -http=:8080 -upstream_url=http://localhost:9000
```

起動後、Manael は指定したアドレスでリクエストを受け付け、アップストリームの画像サーバーへプロキシします。
