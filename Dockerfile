# Start by building the application.
FROM golang:1.26.4-bookworm@sha256:5d2b868674b57c9e48cdd39e891acce4196b6926ca6d11e9c270a8f85106203d AS build

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
