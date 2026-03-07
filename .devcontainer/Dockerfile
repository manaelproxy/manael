FROM mcr.microsoft.com/devcontainers/go:2.0.8-1.25-bookworm@sha256:65af278fcedeef06eb6b672ad4a95c95cf0849cb82d4852ff014719fe6cdf262

# Workaround for devcontainers/images#1797:
# stale Yarn APT source can cause "apt-get update" to fail with NO_PUBKEY.
RUN rm -f /etc/apt/sources.list.d/yarn.list

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends libvips-dev
