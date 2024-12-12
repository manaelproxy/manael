FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm@sha256:2e00578e7c526e76ebd51597260928d8542d27c17ffe6390d79200638c51a3e6

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends libaom-dev libwebp-dev
