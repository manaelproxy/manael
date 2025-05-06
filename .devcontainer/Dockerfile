FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm@sha256:35a12c9f10d6543d24bd3976b0b5f711848deac3b461793b3c410d503990d30e

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends libaom-dev libwebp-dev
