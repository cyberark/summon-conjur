#!/bin/bash

APP="summon-conjur"
WORKDIR="/go/src/github.com/conjurinc/${APP}"

rm -rf pkg

docker run --rm \
-v "$PWD":$WORKDIR \
-w $WORKDIR \
golang:1.6 \
./compile.sh $APP