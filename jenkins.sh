#!/bin/bash -e

./recursive_test.sh
./build.sh

sudo chmod -R 777 pkg/
./package.sh
