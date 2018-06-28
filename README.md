# Manael

Simple WebP proxy server.

## Dependencies

- libjpeg (or libjpeg-turbo)
- libwebp

## Installation

```console
$ go get -u manael.org/x/manael/cmd/manael
```

## Usage

```console
$ manael -http=:8080 -upstream_url=http://localhost:9000
```

## License

[MIT](/LICENSE)
