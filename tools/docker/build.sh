#!/bin/bash

cd `dirname $0`
rm -rf webroot/

cd ../..
make build
cd -

cp -r resource/* webroot/
mkdir webroot/blog
mv webroot/archive \
    webroot/archives.html \
	webroot/categories.html \
	webroot/index.html \
	webroot/blog/

docker build -t levinxo/website:latest .

docker push levinxo/website:latest

rm -rf webroot/

