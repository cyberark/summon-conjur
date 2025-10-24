#!/bin/bash
set -eox pipefail

export PATH="$(pwd):$PATH"

echo "Running Go tests..."

echo "Current dir: $(pwd)"
go test --coverprofile=output/c.out -v ./test/... | tee output/junit.output

go-junit-report < output/junit.output > output/junit.xml

gocov convert output/c.out | gocov-xml > output/coverage.xml

rm output/junit.output
