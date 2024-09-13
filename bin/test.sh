#!/bin/bash -e

export COMPOSE_ARGS='-f docker-compose.yml'

export CONJUR_ACCOUNT=cucumber
export CONJUR_AUTHN_LOGIN=admin

export REGISTRY_URL=${INFRAPOOL_REGISTRY_URL:-"docker.io"}

source ./bin/functions.sh

function finish {
  echo 'Removing environment'
  echo '-----'
  docker compose $COMPOSE_ARGS down -v
}
trap finish EXIT

function main() {
  startConjur
  initEnvironment
  runTests
}

function runTests() {
  local api_key="$(getKeys)"
  
  local service=test

  docker compose $COMPOSE_ARGS build --pull $service

  docker compose $COMPOSE_ARGS run --rm \
    -e GO_TEST_ARGS="$GO_TEST_ARGS" \
    -e CONJUR_AUTHN_API_KEY="$api_key" \
    $service
}

main
