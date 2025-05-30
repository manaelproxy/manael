FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm@sha256:f12d68fb2acd19b5b7a5833ab39e0ff1c16646f698713e3661fce0cb6cb3b590

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends libaom-dev libwebp-dev
