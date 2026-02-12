FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm@sha256:35e445a87924ae04c171d7e52e2a913050dfefbacda3756ef7a307db035704d5

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends libaom-dev libwebp-dev
