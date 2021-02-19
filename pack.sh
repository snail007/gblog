#!/bin/bash
GOVERSION=1.14

cd initialize/
gmct tpl --clean
gmct tpl --dir ../views
gmct static --clean
gmct static --dir ../static
cd ..

rm -rf gblog-release

mkdir gblog-release
mkdir gblog-release/gblog-linux64-release
mkdir gblog-release/gblog-linux32-release
mkdir gblog-release/gblog-win64-release
mkdir gblog-release/gblog-win32-release
mkdir gblog-release/gblog-mac-release
mkdir gblog-release/gblog-arm64-release

cp -R conf gblog-release/gblog-linux64-release
cp -R conf gblog-release/gblog-linux32-release
cp -R conf gblog-release/gblog-win64-release
cp -R conf gblog-release/gblog-win32-release
cp -R conf gblog-release/gblog-mac-release
cp -R conf gblog-release/gblog-arm64-release

set -e

# linux 64
docker run -it --rm \
-v $GOPATH:/go \
-e BUILDDIR=github.com/snail007/gblog \
-e GO111MODULE=on \
-e CGO_ENABLED=1 \
-e GOOS=linux \
-e GOARCH=amd64 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -o gblog-linux64
echo "gblog-linux64 success"

# linux 32
docker run -it --rm \
-v $GOPATH:/go \
-e BUILDDIR=github.com/snail007/gblog \
-e GO111MODULE=on \
-e CGO_ENABLED=1 \
-e GOOS=linux \
-e GOARCH=386 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -o gblog-linux32
echo "gblog-linux32 success"

# windows 64
docker run -it --rm \
-v $GOPATH:/go \
-e BUILDDIR=github.com/snail007/gblog \
-e GO111MODULE=on \
-e CGO_ENABLED=1 \
-e GOOS=windows \
-e GOARCH=386 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -o gblog-win64.exe
echo "gblog-win64.exe success"

# windows 32
docker run -it --rm \
-v $GOPATH:/go \
-e BUILDDIR=github.com/snail007/gblog \
-e GO111MODULE=on \
-e CGO_ENABLED=1 \
-e GOOS=windows \
-e GOARCH=386 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -o gblog-win32.exe
echo "gblog-win32.exe success"

# darwin 64
CGO_ENABLED=1 GO111MODULE=on go build -ldflags "-s -w" -o gblog-mac
echo "gblog-mac success"

# arm64
docker run -it --rm \
-v $GOPATH:/go \
-e BUILDDIR=github.com/snail007/gblog \
-e GO111MODULE=on \
-e CGO_ENABLED=1 \
-e GOOS=linux \
-e GOARCH=arm64 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -o gblog-arm64
echo "gblog-arm64 success"

upx gblog-linux64
upx gblog-linux32
upx gblog-win64.exe
upx gblog-win32.exe
upx gblog-mac
upx gblog-arm64

mv gblog-linux64 gblog-release/gblog-linux64-release/gblog
mv gblog-linux32 gblog-release/gblog-linux32-release/gblog
mv gblog-win64.exe gblog-release/gblog-win64-release/gblog.exe
mv gblog-win32.exe gblog-release/gblog-win32-release/gblog.exe
mv gblog-mac gblog-release/gblog-mac-release/gblog
mv gblog-arm64 gblog-release/gblog-arm64-release/gblog

cd gblog-release
tar zcfv gblog-linux64-release.tar.gz gblog-linux64-release
tar zcfv gblog-linux32-release.tar.gz gblog-linux32-release
tar zcfv gblog-win64-release.tar.gz gblog-win64-release
tar zcfv gblog-win32-release.tar.gz gblog-win32-release
tar zcfv gblog-mac-release.tar.gz gblog-mac-release
tar zcfv gblog-arm64-release.tar.gz gblog-arm64-release
cd ..

cd initialize/
gmct static --clean
gmct tpl --clean
echo "done"


