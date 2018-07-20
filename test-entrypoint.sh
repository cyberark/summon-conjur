#!/bin/bash -eo pipefail

echo "Running tests"
echo "-----"

go clean -i
go install

echo "Running go test with args: $GO_TEST_ARGS"
go test $GO_TEST_ARGS -v "$(go list ./... | grep -v /vendor/)" | tee output/junit.output

go-junit-report < output/junit.output > output/junit.xml

rm output/junit.output
