FROM golang:1.9-stretch

WORKDIR /go/src/app
COPY . /go/src/app

RUN apt-get update \
	&& apt-get install -y --no-install-recommends \
		gcc \
		git \
		libjpeg62-turbo-dev \
		libwebp-dev \
	&& rm -rf /var/lib/apt/lists/* \
	&& go-wrapper download \
	&& go-wrapper install

CMD ["go-wrapper", "run"]
