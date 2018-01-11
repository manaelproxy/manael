FROM golang:1.9.2-alpine3.7 AS build

COPY Gopkg.lock Gopkg.toml *.go /go/src/github.com/ykzts/webp-proxy/

RUN set -ex \
	\
	&& apk add --update --virtual .build-deps \
		gcc \
		git \
		libjpeg-turbo-dev \
		libwebp-dev \
		musl-dev \
	&& go get -u github.com/golang/dep/cmd/dep \
	&& cd /go/src/github.com/ykzts/webp-proxy \
	&& dep ensure \
	&& rm Gopkg.lock Gopkg.toml \
	&& go install \
	&& cd /go \
	&& apk del .build-deps \
	&& rm \
		/var/cache/apk/* \
		/usr/local/bin/* \
		/go/bin/dep \
	&& rm -r \
		/usr/local/go \
		/go/pkg \
		/go/src

FROM alpine:3.7

COPY --from=build /go/bin/webp-proxy /usr/bin/

RUN set -ex \
	\
	&& apk add --update \
		ca-certificates \
		libjpeg-turbo \
		libwebp

CMD ["webp-proxy"]
