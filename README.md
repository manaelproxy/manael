# webp-proxy

Simple WebP proxy server.

## Install

```console
$ go get -u github.com/ykzts/webp-proxy
```

## Usage

```console
$ webp-proxy -port 8080 -upstream-url http://localhost:9000
```

```
# /etc/nginx/conf.d/files.example.com.conf
proxy_cache_path /var/cache/nginx/cache levels=1:2 keys_zone=CACHE:128m max_size=512m inactive=1d;
proxy_temp_path /var/cache/nginx/tmp;

server {
  listen 443 ssl http2;
  server_name files.example.com;
  
  location / {
    add_header X-Cache $upstream_cache_status;
    proxy_pass http://localhost:8080;
    proxy_cache CACHE;
    proxy_cache_valid 200 1d;
    proxy_cache_valid 404 10s;
    proxy_cache_key https://$host$request_uri;
    proxy_cache_lock on;
    proxy_cache_lock_timeout 30s;
  }
}
```

## License

[MIT](/LICENSE)
