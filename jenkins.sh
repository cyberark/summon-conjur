#!/bin/bash -e

./test.sh
./build.sh

chmod -R 777 pkg/
./package.sh
