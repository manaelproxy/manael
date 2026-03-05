# Start by building the application.
FROM golang:1.26.0-bookworm@sha256:2a0ba12e116687098780d3ce700f9ce3cb340783779646aafbabed748fa6677c AS build

ENV LIBWEBP_VERSION=1.2.4
ENV LIBAOM_VERSION=3.8.0
ENV LIBHEIF_VERSION=1.21.2
ENV LIBVIPS_VERSION=8.18.0
ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig

RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		cmake \
		meson \
		ninja-build \
		pkg-config \
		wget \
		yasm \
		libglib2.0-dev \
		libexpat1-dev \
		libjpeg62-turbo-dev \
		libpng-dev \
		libpcre2-dev \
		zlib1g-dev \
	&& rm -rf /var/lib/apt/lists/*

RUN \
	mkdir -p /tmp/src && \
	cd /tmp/src && \
	wget -O libwebp.tar.gz https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-${LIBWEBP_VERSION}.tar.gz && \
	tar -xzf libwebp.tar.gz && \
	rm libwebp.tar.gz && \
	cd /tmp/src/libwebp-${LIBWEBP_VERSION} && \
	./configure --prefix=/usr/local --disable-shared --enable-static && \
	make -j$(nproc) && \
	make install

RUN \
	mkdir -p /tmp/src && \
	cd /tmp/src && \
	wget -O libaom.tar.gz https://storage.googleapis.com/aom-releases/libaom-${LIBAOM_VERSION}.tar.gz && \
	tar -xzf libaom.tar.gz -C /tmp/src && \
	rm libaom.tar.gz && \
	mkdir -p /tmp/src/aom_build && \
	cd /tmp/src/aom_build && \
	cmake /tmp/src/libaom-${LIBAOM_VERSION} \
		-DCMAKE_INSTALL_PREFIX=/usr/local \
		-DBUILD_SHARED_LIBS=OFF && \
	make -j$(nproc) && \
	make install

RUN \
	mkdir -p /tmp/src && \
	cd /tmp/src && \
	wget -O libheif.tar.gz https://github.com/strukturag/libheif/releases/download/v${LIBHEIF_VERSION}/libheif-${LIBHEIF_VERSION}.tar.gz && \
	tar -xzf libheif.tar.gz && \
	rm libheif.tar.gz && \
	mkdir -p /tmp/src/libheif-${LIBHEIF_VERSION}/build && \
	cd /tmp/src/libheif-${LIBHEIF_VERSION}/build && \
	cmake .. \
		-DCMAKE_INSTALL_PREFIX=/usr/local \
		-DCMAKE_PREFIX_PATH=/usr/local \
		-DBUILD_SHARED_LIBS=OFF \
		-DWITH_AOM_ENCODER=ON \
		-DWITH_AOM_DECODER=ON && \
	make -j$(nproc) && \
	make install

RUN \
	mkdir -p /tmp/src && \
	cd /tmp/src && \
	wget -O libvips.tar.xz https://github.com/libvips/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.xz && \
	tar -xJf libvips.tar.xz && \
	rm libvips.tar.xz && \
	cd /tmp/src/vips-${LIBVIPS_VERSION} && \
	meson setup build \
		--prefix=/usr/local \
		--default-library=static \
		-Dwebp=enabled \
		-Dheif=enabled \
		-Djpeg=enabled \
		-Dpng=enabled \
		-Dmagick=disabled \
		-Dcgif=disabled \
		-Djpeg-xl=disabled \
		-Dopenslide=disabled \
		-Dtiff=disabled \
		-Dpdfium=disabled \
		-Dmodules=disabled && \
	cd /tmp/src/vips-${LIBVIPS_VERSION}/build && \
	ninja && \
	ninja install

WORKDIR /go/src/manael
COPY . .

RUN go mod download
RUN CGO_LDFLAGS="$(pkg-config --static --libs vips)" \
	go build -ldflags '-extldflags "-static"' -o /go/bin/manael ./cmd/manael

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian12@sha256:937c7eaaf6f3f2d38a1f8c4aeff326f0c56e4593ea152e9e8f74d976dde52f56
COPY --from=build /go/bin/manael /
CMD ["/manael"]
