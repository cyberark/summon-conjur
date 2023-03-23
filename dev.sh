#!/usr/bin/env bash -x

# function finish {
#   echo 'Removing environment'
#   echo '-----'
#   docker-compose down -v
# }
# trap finish EXIT
#

export CONJUR_ACCOUNT=cucumber
export CONJUR_AUTHN_LOGIN=admin

source $(dirname $0)/bin/functions.sh

function main() {
  startConjur 'all'
  initEnvironment 'all'
  runDevelopment
}

function runDevelopment() {
  local keys=( $(getKeys) )
  local api_key=${keys[0]}

  export CONJUR_AUTHN_API_KEY="$api_key"
  docker-compose up -d cli

  docker-compose build --pull dev

  docker-compose run -d \
    --service-ports \
    dev
}

main
