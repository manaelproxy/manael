# Start by building the application.
FROM golang:1.26.2-bookworm@sha256:e1b367ff64fba58a08544465ad6574a619be7517f9a517ff982374d1fd1085fa AS build

RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		libvips-dev \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/manael
COPY . .

RUN go mod download
RUN go build -o /go/bin/manael ./cmd/manael

# Now copy it into our base image.
FROM debian:bookworm-slim@sha256:4724b8cc51e33e398f0e2e15e18d5ec2851ff0c2280647e1310bc1642182655d
RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		libvips42 \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=build /go/bin/manael /
CMD ["/manael"]
