#!/bin/sh -x
H=$(readlink -f "$0" | md5sum | head -c 8)
GOCACHE=go-cache-$H
GOPATH=go-path-$H
docker volume create "$GOCACHE"
docker volume create "$GOPATH"
docker run -v "$GOCACHE:/.cache" -v "$GOPATH:/go" alpine sh -ec "mkdir -p /.cache/go-build /go/bin /go/pkg /go/src; chown $(id -u):$(id -g) /.cache/go-build /go/bin /go/pkg /go/src"
docker run -v "$GOCACHE:/.cache" -v "$GOPATH:/go" -v "$PWD:/srv" -w /srv -u "$(id -u):$(id -g)" golang go "$@"
