#!/bin/bash -e

export CONJUR_ACCOUNT=cucumber
export CONJUR_AUTHN_LOGIN=admin

source functions.sh

function finish {
  echo 'Removing environment'
  echo '-----'
  docker-compose down -v
}
trap finish EXIT

function main() {
  startConjur
  initEnvironment
  runTests
}

function runTests() {
  local keys=( $(getKeys) )
  local api_key=${keys[0]}
  local api_key_v4=${keys[1]}
  local ssl_cert_v4="$(getCert)"
  local service=test
  
  docker-compose build --pull $service
  
  docker-compose run --rm \
    -e CONJUR_AUTHN_API_KEY="$api_key" \
    -e CONJUR_V4_AUTHN_API_KEY="$api_key_v4" \
    -e CONJUR_V4_SSL_CERTIFICATE="$ssl_cert_v4" \
    $service
}

main
