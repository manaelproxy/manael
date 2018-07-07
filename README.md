# Manael

[![GoDoc](https://godoc.org/manael.org/x/manael?status.svg)](https://godoc.org/manael.org/x/manael) [![CircleCI](https://circleci.com/gh/manaelproxy/manael/tree/master.svg?style=shield)](https://circleci.com/gh/manaelproxy/manael/tree/master)

A simple HTTP proxy for processing images.

## Dependencies

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
