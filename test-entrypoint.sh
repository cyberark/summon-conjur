#!/bin/bash -eo pipefail

echo "Running tests"
echo "-----"

GO_TEST_ARGS='--tags=all -vet=off'

echo "Running go test with args: $GO_TEST_ARGS"
go test $GO_TEST_ARGS -v | tee output/junit.output

go-junit-report < output/junit.output > output/junit.xml

rm output/junit.output
