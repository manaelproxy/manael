---
title: Getting Started
weight: 10
---

Start Manael by specifying the address to listen on and the upstream server URL:

```console
manael -http=:8080 -upstream_url=http://localhost:9000
```

After startup, Manael listens on the configured address and proxies image requests to the upstream server.
