#!/bin/bash

cd initialize/
gmct tpl --clean
gmct tpl --dir ../views
gmct static --clean
gmct static --dir ../static
cd ..

#go mod vendor

rm -rf gblog-release
mkdir gblog-release
cp -R conf gblog-release/

cd gblog-release/
CGO_ENABLED=1 xgo --targets=linux/arm-5,linux/arm-6,linux/arm-7,linux/arm64,windows/amd64,windows/386,linux/amd64,linux/386,darwin-10.10/amd64,darwin-10.10/386 -go latest -ldflags "-s -w" ../
#CGO_ENABLED=1 xgo.karalabe --targets=linux/amd64 -go latest ../
#CGO_ENABLED=1 go build -ldflags "-s -w" -o gblog
cd ../

rm -rf gblog-release.tar.gz
tar zcfv gblog-release.tar.gz gblog-release

cd initialize/
gmct static --clean
gmct tpl --clean
cd ..

#rm -rf vendor


