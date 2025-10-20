#!/bin/bash -e

export CONJUR_ACCOUNT=cucumber
export CONJUR_AUTHN_LOGIN=admin

export REGISTRY_URL=${INFRAPOOL_REGISTRY_URL:-"docker.io"}

source ./bin/functions.sh

function finish {
  echo 'Removing environment'
  echo '-----'
  docker compose down -v
}
trap finish EXIT

function failed {
  echo 'TESTS FAILED'
  echo '-----'
  echo 'Conjur logs:'
  docker compose logs conjur || true
  echo '-----'
  exit 1
}

function main() {
  startConjur
  initEnvironment
  prepareGCP
  runTests
}

function prepareGCP() {
  if [[ "$INFRAPOOL_TEST_GCP" == "true" ]]; then
    GCP_PROJECT_ID=""
    GCP_ID_TOKEN=""
    if [[ -f "gcp/project-id" ]]; then
      read -r GCP_PROJECT_ID < "gcp/project-id"
    fi
    if [[ -f "gcp/token" ]]; then
      read -r GCP_ID_TOKEN < "gcp/token"
    fi
    if [[ -z "$GCP_PROJECT_ID" || -z "$GCP_ID_TOKEN" ]]; then
      echo "GCP_PROJECT_ID and GCP_ID_TOKEN must be set to run GCP tests"
      failed
    fi
    export GCP_PROJECT_ID
    export GCP_ID_TOKEN
  fi
}

function runTests() {
  local api_key="$(getKeys)"
  
  local service=test

  docker compose build --pull $service

  docker compose run --rm \
    -e GO_TEST_ARGS="$GO_TEST_ARGS" \
    -e CONJUR_AUTHN_API_KEY="$api_key" \
    -e TEST_AWS="$INFRAPOOL_TEST_AWS" \
    -e TEST_GCP="$INFRAPOOL_TEST_GCP" \
    -e GCP_PROJECT_ID \
    -e GCP_ID_TOKEN \
    -e TEST_AZURE="$INFRAPOOL_TEST_AZURE" \
    -e AZURE_SUBSCRIPTION_ID \
    -e AZURE_RESOURCE_GROUP \
    $service || failed
}

main
