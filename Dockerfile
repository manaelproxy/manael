# Start by building the application.
FROM golang:1.26.1-bookworm@sha256:ab3d6955bbc813a0f3fdf220c1d817dd89c0b3f283777db8ece4a32fe7858edd AS build

RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		libvips-dev \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/manael
COPY . .

RUN go mod download
RUN go build -o /go/bin/manael ./cmd/manael

# Now copy it into our base image.
FROM debian:bookworm-slim@sha256:74d56e3931e0d5a1dd51f8c8a2466d21de84a271cd3b5a733b803aa91abf4421
RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		libvips42 \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=build /go/bin/manael /
CMD ["/manael"]
