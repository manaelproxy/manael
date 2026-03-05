# Start by building the application.
FROM golang:1.26.0-bookworm@sha256:2a0ba12e116687098780d3ce700f9ce3cb340783779646aafbabed748fa6677c AS build

RUN apt-get update && \
	apt-get install -y --no-install-recommends libvips-dev pkg-config && \
	rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/manael
COPY . .

RUN go mod download
RUN go build -o /go/bin/manael ./cmd/manael

# Now copy it into our base image.
FROM debian:bookworm-slim

RUN apt-get update && \
	apt-get install -y --no-install-recommends libvips42 && \
	rm -rf /var/lib/apt/lists/*

COPY --from=build /go/bin/manael /
CMD ["/manael"]
