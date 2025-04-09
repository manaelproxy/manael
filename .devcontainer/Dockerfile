FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm@sha256:ee28302232bca53c6cfacf0b00a427ebbda10b33731c78d3dcf9f59251b23c9c

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends libaom-dev libwebp-dev
