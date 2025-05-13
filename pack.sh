#!/bin/bash -x
GOVERSION=1.16

go mod tidy

cd initialize/
gmct tpl --clean
gmct tpl --dir ../views
gmct static --clean
gmct static --dir ../static
cd ..

rm -rf gblog-release
rm -rf conf_release

mkdir gblog-release
mkdir gblog-release/gblog-linux64-release
mkdir gblog-release/gblog-linux32-release
mkdir gblog-release/gblog-win64-release
mkdir gblog-release/gblog-win32-release
mkdir gblog-release/gblog-mac-release
mkdir gblog-release/gblog-arm64-release
#with bleve
mkdir gblog-release/gblog-linux64-release-bleve
mkdir gblog-release/gblog-linux32-release-bleve
mkdir gblog-release/gblog-win64-release-bleve
mkdir gblog-release/gblog-win32-release-bleve
mkdir gblog-release/gblog-mac-release-bleve

set -e

# linux 64
docker run -it --rm \
-v $GOPATH:/go \
-w /go/src/github.com/snail007/gblog \
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
-w /go/src/github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
-e CGO_ENABLED=1 \
-e GOOS=linux \
-e GOARCH=amd64 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -tags "bleve" -o gblog-linux64-bleve
echo "gblog-linux64-bleve success"

# linux 32
docker run -it --rm \
-v $GOPATH:/go \
-w /go/src/github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
-e CGO_ENABLED=1 \
-e GOOS=linux \
-e GOARCH=386 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -o gblog-linux32
echo "gblog-linux32 success"

# linux 32 with bleve
docker run -it --rm \
-v $GOPATH:/go \
-w /go/src/github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
-e CGO_ENABLED=1 \
-e GOOS=linux \
-e GOARCH=386 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -tags "bleve" -o gblog-linux32-bleve
echo "gblog-linux32-bleve success"


# windows 64
docker run -it --rm \
-v $GOPATH:/go \
-w /go/src/github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
-e CGO_ENABLED=1 \
-e GOOS=windows \
-e GOARCH=386 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -o gblog-win64.exe
echo "gblog-win64.exe success"

# windows 64 with bleve
docker run -it --rm \
-v $GOPATH:/go \
-w /go/src/github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
-e CGO_ENABLED=1 \
-e GOOS=windows \
-e GOARCH=386 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -tags "bleve" -o gblog-win64-bleve.exe
echo "gblog-win64-bleve.exe success"


# windows 32
docker run -it --rm \
-v $GOPATH:/go \
-w /go/src/github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
-e CGO_ENABLED=1 \
-e GOOS=windows \
-e GOARCH=386 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -o gblog-win32.exe
echo "gblog-win32.exe success"

# windows 32 with bleve
docker run -it --rm \
-v $GOPATH:/go \
-w /go/src/github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
-e CGO_ENABLED=1 \
-e GOOS=windows \
-e GOARCH=386 \
snail007/golang:$GOVERSION \
go build -ldflags "-s -w" -tags "bleve" -o gblog-win32-bleve.exe
echo "gblog-win32-bleve.exe success"


# darwin 64
CGO_ENABLED=1 GO111MODULE=on go build -ldflags "-s -w" -o gblog-mac
echo "gblog-mac success"

# darwin 64 with bleve
CGO_ENABLED=1 GO111MODULE=on go build -ldflags "-s -w" -tags "bleve" -o gblog-mac-bleve
echo "gblog-mac-bleve success"

# arm64
docker run -it --rm \
-v $GOPATH:/go \
-w /go/src/github.com/snail007/gblog \
-e GO111MODULE=on \
-e GOSUMDB=off \
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
upx gblog-arm64
#upx gblog-mac
# with bleve
upx gblog-linux64-bleve
upx gblog-linux32-bleve
upx gblog-win64-bleve.exe
upx gblog-win32-bleve.exe
#upx gblog-mac-bleve

mv gblog-linux64 gblog-release/gblog-linux64-release/gblog
mv gblog-linux32 gblog-release/gblog-linux32-release/gblog
mv gblog-win64.exe gblog-release/gblog-win64-release/gblog.exe
mv gblog-win32.exe gblog-release/gblog-win32-release/gblog.exe
mv gblog-mac gblog-release/gblog-mac-release/gblog
mv gblog-arm64 gblog-release/gblog-arm64-release/gblog
# with bleve
mv gblog-linux64-bleve gblog-release/gblog-linux64-release-bleve/gblog
mv gblog-linux32-bleve gblog-release/gblog-linux32-release-bleve/gblog
mv gblog-win64-bleve.exe gblog-release/gblog-win64-release-bleve/gblog.exe
mv gblog-win32-bleve.exe gblog-release/gblog-win32-release-bleve/gblog.exe
mv gblog-mac-bleve gblog-release/gblog-mac-release-bleve/gblog

cd gblog-release
tar zcfv gblog-linux64-release.tar.gz gblog-linux64-release
tar zcfv gblog-linux32-release.tar.gz gblog-linux32-release
tar zcfv gblog-win64-release.tar.gz gblog-win64-release
tar zcfv gblog-win32-release.tar.gz gblog-win32-release
tar zcfv gblog-mac-release.tar.gz gblog-mac-release
tar zcfv gblog-arm64-release.tar.gz gblog-arm64-release
# with bleve
tar zcfv gblog-linux64-release-bleve.tar.gz gblog-linux64-release-bleve
tar zcfv gblog-linux32-release-bleve.tar.gz gblog-linux32-release-bleve
tar zcfv gblog-win64-release-bleve.tar.gz gblog-win64-release-bleve
tar zcfv gblog-win32-release-bleve.tar.gz gblog-win32-release-bleve
tar zcfv gblog-mac-release-bleve.tar.gz gblog-mac-release-bleve
cd ..

cd initialize/
gmct static --clean
gmct tpl --clean
echo "done"


