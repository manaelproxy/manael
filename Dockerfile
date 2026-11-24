# Start by building the application.
FROM golang:1.27.4-bookworm@sha256:b754301872c53d3dff666ec39e3fae971dc6f5ecd5a9b01fdefabbebb0f59901 AS build

RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		libvips-dev \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/manael
COPY . .

RUN go mod download
RUN go build -o /go/bin/manael ./cmd/manael

# Now copy it into our base image.
FROM debian:bookworm-slim@sha256:e67eeead736ff7bbe258ef088bff683635c561b1fcfefad0806c1df4ab49ce0e
RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		libvips42 \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=build /go/bin/manael /
CMD ["/manael"]
