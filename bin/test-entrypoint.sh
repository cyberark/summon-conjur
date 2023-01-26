#!/bin/bash -eo pipefail

export PATH="$(pwd):$PATH"
echo "Path: $PATH"

echo "Running tests..."

TEST_PARAMS="-run TestPackage*"

echo "Running go tests: $TEST_PARAMS"
echo "Current dir: $(pwd)"

set -x
go test --coverprofile=output/c.out -v ./test/... $TEST_PARAMS | tee output/junit.output

go-junit-report < output/junit.output > output/junit.xml

gocov convert output/c.out | gocov-xml > output/coverage.xml

rm output/junit.output
