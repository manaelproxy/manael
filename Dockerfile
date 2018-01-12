FROM golang:1.9.2-alpine3.7 AS build

COPY Gopkg.lock Gopkg.toml *.go /go/src/github.com/ykzts/webp-proxy/

RUN set -ex \
	\
	&& echo http://dl-cdn.alpinelinux.org/alpine/v3.7/main > /etc/apk/repositories \
	&& apk --update-cache add --virtual .build-deps \
		gcc \
		git \
		libjpeg-turbo-dev \
		libwebp-dev \
		musl-dev \
	\
	&& go get -u github.com/golang/dep/cmd/dep \
	&& cd /go/src/github.com/ykzts/webp-proxy \
	&& dep ensure \
	&& go install -ldflags '-s -w' \
	\
	&& cd /tmp \
	&& mkdir -p rootfs/usr/bin \
	&& cp /go/bin/webp-proxy rootfs/usr/bin \
	&& mkdir -p rootfs/lib rootfs/usr/lib \
	&& for lib in $(ldd rootfs/usr/bin/webp-proxy | awk '{ print $(NF-1) }'); do \
		cp $lib rootfs$lib; \
	done \
	&& mkdir -p rootfs/etc/ssl/certs \
	&& cp /etc/ssl/certs/ca-certificates.crt rootfs/etc/ssl/certs \
	\
	&& apk del .build-deps \
	&& rm \
		/var/cache/apk/* \
		/usr/local/bin/* \
		/go/bin/dep \
	&& rm -r \
		/usr/local/go \
		/go/pkg \
		/go/src

FROM scratch

COPY --from=build /tmp/rootfs /

CMD ["webp-proxy"]
