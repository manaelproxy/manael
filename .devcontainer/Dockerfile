FROM mcr.microsoft.com/devcontainers/go:2.1.1-1.25-bookworm@sha256:2765adcd37710f2e24c67837cf9c88ae1d88708f405c15a88b307711268db63e

# Workaround for devcontainers/images#1797:
# stale Yarn APT source can cause "apt-get update" to fail with NO_PUBKEY.
RUN rm -f /etc/apt/sources.list.d/yarn.list

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends libvips-dev
