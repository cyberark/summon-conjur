#!/bin/bash -e

CONJUR_TYPE="${1:-all}"  # Type of Conjur to test against, 'all', 'oss' or 'enterprise'

if [[ "$CONJUR_TYPE" == "all" || "$CONJUR_TYPE" == "enterprise" ]]; then
  export COMPOSE_ARGS='-f docker-compose.yml -f docker-compose.enterprise.yml'
fi

export CONJUR_ACCOUNT=cucumber
export CONJUR_AUTHN_LOGIN=admin

source ./bin/functions.sh

function finish {
  echo 'Removing environment'
  echo '-----'
  docker-compose $COMPOSE_ARGS down -v
}
trap finish EXIT

function main() {
  startConjur $CONJUR_TYPE
  initEnvironment $CONJUR_TYPE
  runTests $CONJUR_TYPE
}

function runTests() {
  local conjurType="$1"

  if [[ "$conjurType" == "all" || "$conjurType" == "enterprise" ]]; then
    local api_key_v4="$(getKeys 'enterprise')"
    local ssl_cert_v4="$(getCert)"
  fi
  if [[ "$conjurType" == "all" || "$conjurType" == "oss" ]]; then
    local api_key="$(getKeys 'oss')"
  fi

  local service=test

  docker-compose $COMPOSE_ARGS build --pull $service

  docker-compose $COMPOSE_ARGS run --rm \
    -e CONJUR_TYPE="$CONJUR_TYPE" \
    -e GO_TEST_ARGS="$GO_TEST_ARGS" \
    -e CONJUR_AUTHN_API_KEY="$api_key" \
    -e CONJUR_V4_AUTHN_API_KEY="$api_key_v4" \
    -e CONJUR_V4_SSL_CERTIFICATE="$ssl_cert_v4" \
    $service
}

main
