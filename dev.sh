#!/usr/bin/env bash

function finish {
  echo 'Removing environment'
  echo '-----'
  docker-compose down -v
}
trap finish EXIT

function main() {
  startConjur
  runDevelopment
}

function startConjur() {
  # Start a local Conjur to test against
  docker-compose pull postgres conjur
  docker-compose up -d conjur
}

function runDevelopment() {
  docker-compose build --pull tester
  # Delay to allow time for Conjur to come up
  # TODO: remove this once we have HEALTHCHECK in place
  docker-compose run --rm tester ./wait_for_server.sh

  # Run development
  local api_key=$(docker-compose exec -T conjur rails r "print Credentials['cucumber:user:admin'].api_key")

  docker-compose run --rm \
    --entrypoint "bash -c './convey.sh& bash'" \
    --service-ports \
    -e CONJUR_AUTHN_API_KEY="$api_key" \
    tester
}

main
