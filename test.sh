#!/bin/bash -e

function finish {
  echo 'Removing environment'
  echo '-----'
  docker-compose down -v
}
trap finish EXIT

function main() {
  startConjur
  runTests
}

function startConjur() {
  # Start a local Conjur to test against
  docker-compose pull postgres conjur
  docker-compose up -d conjur
}

function runTests() {
  docker-compose build --pull tester
  # Delay to allow time for Conjur to come up
  # TODO: remove this once we have HEALTHCHECK in place
  docker-compose run --rm tester ./wait_for_server.sh

  # Execute tests
  local api_key=$(docker-compose exec -T conjur rails r "print Credentials['cucumber:user:admin'].api_key")

  docker-compose run --rm \
    -e CONJUR_API_KEY="$api_key" \
    tester
}

main
