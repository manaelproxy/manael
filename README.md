# Manael

[![GoDoc](https://godoc.org/manael.org/x/manael?status.svg)](https://godoc.org/manael.org/x/manael)
[![Go](https://github.com/manaelproxy/manael/workflows/Go/badge.svg)](https://github.com/manaelproxy/manael/actions?query=workflow%3AGo)
[![Codecov](https://codecov.io/gh/manaelproxy/manael/branch/main/graph/badge.svg)](https://codecov.io/gh/manaelproxy/manael)

Manael is a simple HTTP proxy for processing images.

## Installation

- [Download latest binary](https://github.com/manaelproxy/manael/releases/latest)

### Build from source

Building from source requires [libvips](https://www.libvips.org/) development headers in addition to Go and Git.

On Debian/Ubuntu:

```console
sudo apt-get install -y libvips-dev
```

On macOS (Homebrew):

```console
brew install vips
```

## Usage

Start the proxy server:

```console
manael -http=:8080 -upstream_url=http://localhost:9000
```

To convert a JPEG image to WebP, send a request with an `Accept: image/webp` header. Manael will automatically convert the image if the upstream server returns a JPEG or PNG:

```console
curl -sI -H "Accept: image/webp" http://localhost:8080/image.jpg
```

The response will have `Content-Type: image/webp` when conversion succeeds.

## License

[MIT](/LICENSE)
