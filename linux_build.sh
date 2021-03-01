#!/bin/bash
GOVERSION=1.16

go mod tidy

cd initialize/
gmct tpl --clean
gmct tpl --dir ../views
gmct static --clean
gmct static --dir ../static
cd ..

set -e

# linux 64
docker run -it --rm \
-v $GOPATH:/go \
-e BUILDDIR=github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
-e CGO_ENABLED=1 \
-e GOOS=linux \
-e GOARCH=amd64 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -o gblog-linux64
echo "gblog-linux64 success"

# linux 64 with bleve
docker run -it --rm \
-v $GOPATH:/go \
-e BUILDDIR=github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
-e CGO_ENABLED=1 \
-e GOOS=linux \
-e GOARCH=amd64 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -tags "bleve" -o gblog-linux64-bleve
echo "gblog-linux64-bleve success"

cd initialize/
gmct static --clean
gmct tpl --clean
echo "done"


