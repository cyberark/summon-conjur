#!/bin/bash -eo pipefail

export PATH="$(pwd):$PATH"
echo "Path: $PATH"

echo "Running tests..."

# Type of Conjur to test against, 'all', 'oss' or 'enterprise'
CONJUR_TYPE="${CONJUR_TYPE:-all}"
echo "Test coverage: $CONJUR_TYPE"

TEST_PARAMS="-run TestPackage*"

if [[ "${CONJUR_TYPE}" == "enterprise" ]]; then
  TEST_PARAMS="-run TestPackageEnterprise"
fi

if [[ "${CONJUR_TYPE}" == "oss" ]]; then
  TEST_PARAMS="-run TestPackageOSS"
fi

echo "Running go tests: $TEST_PARAMS"
echo "Current dir: $(pwd)"

set -x
go test -v ./test/... $TEST_PARAMS | tee output/junit.output

go-junit-report < output/junit.output > output/junit.xml

rm output/junit.output
