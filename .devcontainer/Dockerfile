FROM mcr.microsoft.com/devcontainers/go:2.1.5-1.25-bookworm@sha256:f8f841c474489d6cd5b67cb30ee0e56ae5f1379f0bdf2679191573c41b5e11f8

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
