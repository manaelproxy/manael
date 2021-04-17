# Installation Guide

## Using the Docker {#using-the-docker}

The Manael recommends running on [Docker](https://www.docker.com/). A Docker image is published on [Docker Hub](https://hub.docker.com/).

Get it with `docker pull manael/manael:latest` command before running the Manael with the Docker. Using the Docker eliminates a need to add unnecessary files to your environment.

## Using a built binary {#using-a-built-binary}

You can download the Manael built for 64bit GNU/Linux.

### 1. Create a working directory {#1-create-a-working-directory}

First, create a working directory to extract the downloaded file when installing Manael.

```console
$ mkdir manael
$ cd manael
```

### 2. Download {#2-download}

Download the latest version of Manael (`manael_1.x.y_Linux_x86_64.tar.gz`) from [release page](https://github.com/manaelproxy/manael/releases). Then, extract the downloaded file to the directory created in 1.

```console
$ wget https://github.com/manaelproxy/manael/releases/download/v1.x.y/manael_1.x.y_Linux_x86_64.tar.gz
$ tar xf manael_1.x.y_Linux_x86_64.tar.gz
```

### 3. Install {#3-install}

Use the `install` command to copy the file. You can do the same thing with the `cp` and `mv` commands, but using the `install` command gives the executable the appropriate execution permissions.

```console
$ sudo install manael /usr/local/bin
```

## Build from a source code {#build-from-a-source-code}

A source code is hosted on [GitHub](https://github.com/manaelproxy/manael). The Manael is written in [Go](https://golang.org/).

```console
$ go build -o manael cmd/manael/main.go
```
