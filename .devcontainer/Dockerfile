FROM mcr.microsoft.com/devcontainers/go:1-1.21-bookworm

ENV LIBWEBP_VERSION 1.3.2

RUN \
	apt-get update && \
	export DEBIAN_FRONTEND=noninteractive && \
	apt-get install -y libaom-dev && \
	\
	mkdir -p /tmp/src && \
	cd /tmp/src && \
	wget -O libwebp.tar.gz https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-${LIBWEBP_VERSION}.tar.gz && \
	tar -xzf libwebp.tar.gz -C /tmp/src && \
	rm libwebp.tar.gz && \
	cd /tmp/src/libwebp-${LIBWEBP_VERSION} && \
	./configure && \
	make -j4 && \
	make install
