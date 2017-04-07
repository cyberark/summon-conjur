#!/bin/bash

WORKDIR=$PWD
for i in $(find . -name 'test.sh');
do
   cd $(dirname $i)
   echo "running tests for: $PWD"
   ./test.sh
   cd $WORKDIR
done