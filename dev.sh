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
  docker-compose pull postgres conjur cuke-master
  docker-compose up -d conjur cuke-master
}

function runDevelopment() {
  docker-compose build --pull tester
  # Delay to allow time for Conjur to come up
  # TODO: remove this once we have HEALTHCHECK in place
  docker-compose run --rm tester ./wait_for_server.sh

  # Run development
  local api_key=$(docker-compose exec -T conjur rails r "print Credentials['cucumber:user:admin'].api_key")

  docker-compose exec -T cuke-master bash -c "conjur authn login -u admin -p secret"
  docker-compose exec -T cuke-master bash -c "conjur variable create existent-variable-with-undefined-value"
  docker-compose exec -T cuke-master bash -c "conjur variable create existent-variable-with-defined-value"
  docker-compose exec -T cuke-master bash -c "conjur variable values add existent-variable-with-defined-value existent-variable-defined-value"

  local api_key_v4=$(docker-compose exec -T cuke-master bash -c "conjur user rotate_api_key")
  local ssl_cert_v4=$(docker-compose exec -T cuke-master bash -c "cat /opt/conjur/etc/ssl/ca.pem")

  docker-compose run --rm \
    --entrypoint "bash -c './convey.sh& bash'" \
    --service-ports \
    -e CONJUR_AUTHN_API_KEY="$api_key" \
    -e CONJUR_V4_AUTHN_API_KEY="$api_key_v4" \
    -e CONJUR_V4_SSL_CERTIFICATE="$ssl_cert_v4" \
    tester
}

main
