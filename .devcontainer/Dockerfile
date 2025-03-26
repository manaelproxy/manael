FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm@sha256:eaea1ec5fbab39d9c7bd26a1792d44519f9c458808736f76c1f86d74c7727c54

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends libaom-dev libwebp-dev
