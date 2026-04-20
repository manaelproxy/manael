FROM mcr.microsoft.com/devcontainers/go:2.1.2-1.25-bookworm@sha256:80694b9e519d41822772c5e4b2f7903343f7591ca5db064613ec7e53b855dc31

# renovate: datasource=github-releases depName=gohugoio/hugo versioning=semver-coerced
ARG HUGO_VERSION=0.160.1

# Workaround for devcontainers/images#1797:
# stale Yarn APT source can cause "apt-get update" to fail with NO_PUBKEY.
RUN rm -f /etc/apt/sources.list.d/yarn.list

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends libvips-dev

RUN set -eux; \
	arch="$(dpkg --print-architecture)"; \
	curl -fsSL https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-${arch}.tar.gz -o /tmp/hugo.tgz; \
	tar -xzf /tmp/hugo.tgz -C /tmp hugo; \
	install -m 0755 /tmp/hugo /usr/local/bin/hugo; \
	rm -f /tmp/hugo /tmp/hugo.tgz
