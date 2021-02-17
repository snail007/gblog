#!/bin/bash

cd initialize/
gmct tpl --clean
gmct tpl --dir ../views
gmct static --clean
gmct static --dir ../static
cd ..

go build -ldflags "-s -w" -o gblog

rm -rf gblog-release
mkdir gblog-release

mv gblog gblog-release/
cp -R conf gblog-release/

rm -rf gblog-release.tar.gz

tar zcfv gblog-release.tar.gz gblog-release

cd initialize/
gmct tpl --clean
gmct static --clean
cd ..

rm -rf admin-release