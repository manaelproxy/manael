# Installation Guide

## Using Docker {#using-the-docker}

It is recommended to run Manael with [Docker](https://www.docker.com/). A Docker image for running Manael is published on [Docker Hub](https://hub.docker.com/).

Get the image with `docker pull manael/manael:latest` command before running Manael with Docker. Using Docker eliminates a need to add unnecessary files to your environment.

## Using a binary {#using-a-built-binary}

You can download the Manael build for 64bit GNU/Linux.

### 1. Create a working directory {#1-create-a-working-directory}

First, create a working directory to extract the downloaded file when installing Manael.

```console
$ mkdir manael
$ cd manael
```

### 2. Download {#2-download}

Download the latest version of Manael (`manael_1.x.y_Linux_x86_64.tar.gz`) from the [release page](https://github.com/manaelproxy/manael/releases) on GitHub. Then, extract the downloaded file to the directory created in step 1.

```console
$ wget https://github.com/manaelproxy/manael/releases/download/v1.x.y/manael_1.x.y_Linux_x86_64.tar.gz
$ tar xf manael_1.x.y_Linux_x86_64.tar.gz
```

### 3. Install {#3-install}

Use the `install` command to copy the file. You can do the same thing with the `cp` and `mv` commands, but using the `install` command gives the executable the appropriate execution permissions.

```console
$ sudo install manael /usr/local/bin
```

## Build from source {#build-from-a-source-code}

The source code is hosted on [GitHub](https://github.com/manaelproxy/manael), and Manael is written in [Go](https://go.dev/). To install Manael, make sure to install Go and [Git](https://git-scm.com/) first, and [copy the repository](https://gist.github.com/natedana/cc71d496b611e70673cab5e8f5a78485).

```console
$ go build -o manael cmd/manael/main.go
```
