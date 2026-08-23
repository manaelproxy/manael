# Start by building the application.
FROM golang:1.27.0-bookworm@sha256:484ef6066fa69acb059fdfeda7ba2b8f7391f2ef6abc6f9b8411e669ebd56466 AS build

RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		libvips-dev \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/manael
COPY . .

RUN go mod download
RUN go build -o /go/bin/manael ./cmd/manael

# Now copy it into our base image.
FROM debian:bookworm-slim@sha256:60eac759739651111db372c07be67863818726f754804b8707c90979bda511df
RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		libvips42 \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=build /go/bin/manael /
CMD ["/manael"]
