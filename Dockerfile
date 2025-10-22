# Start by building the application.
FROM golang:1.25.3-bookworm@sha256:ee420c17fa013f71eca6b35c3547b854c838d4f26056a34eb6171bba5bf8ece4 AS build

ENV LIBAOM_VERSION=3.8.0
ENV LIBWEBP_VERSION=1.2.4

RUN \
	apt-get update && \
	apt-get install -y cmake yasm && \
	\
	mkdir -p /tmp/src && \
	cd /tmp/src && \
	wget -O libwebp.tar.gz https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-${LIBWEBP_VERSION}.tar.gz && \
	tar -xzf libwebp.tar.gz -C /tmp/src && \
	rm libwebp.tar.gz && \
	cd /tmp/src/libwebp-${LIBWEBP_VERSION} && \
	./configure --prefix /tmp/libwebp && \
	make -j4 && \
	make install && \
	\
	mkdir -p /tmp/src && \
	cd /tmp/src && \
	wget -O libaom.tar.gz https://storage.googleapis.com/aom-releases/libaom-${LIBAOM_VERSION}.tar.gz && \
	mkdir -p /tmp/src/libaom-${LIBAOM_VERSION} && \
	tar -xzf libaom.tar.gz -C /tmp/src && \
	rm libaom.tar.gz && \
	mkdir -p /tmp/src/aom_build && \
	cd /tmp/src/aom_build && \
	cmake /tmp/src/libaom-${LIBAOM_VERSION} -DCMAKE_INSTALL_PREFIX=/tmp/libaom && \
	make -j4 && \
	make install

WORKDIR /go/src/manael
COPY . .

ENV CGO_CFLAGS="-I/tmp/libwebp/include -I/tmp/libaom/include"
ENV CGO_LDFLAGS="-L/tmp/libwebp/lib -lwebp -L/tmp/libaom/lib -laom -lm"

RUN go mod download
RUN go build -ldflags '-extldflags "-static"' -o /go/bin/manael ./cmd/manael

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian12@sha256:1951bedd9ab20dd71a5ab11b3f5a624863d7af4109f299d62289928b9e311d5d
COPY --from=build /go/bin/manael /
CMD ["/manael"]
