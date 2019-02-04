function startConjur() {
  local conjurType="$1"
  local services='conjur cuke-master'

  if [[ "$conjurType" == "oss" ]]; then
    services='conjur'
  elif [[ "$conjurType" == "enterprise" ]]; then
    services='cuke-master'
  fi

  docker-compose $COMPOSE_ARGS pull $services
  docker-compose $COMPOSE_ARGS up -d $services
}

exec_on() {
  local container=$1; shift
  docker exec $(docker-compose $COMPOSE_ARGS ps -q $container) "$@"
}

function initEnvironment() {
  local conjurType="$1"

  if [[ "$conjurType" == "all" || "$conjurType" == "oss" ]]; then
    exec_on conjur conjurctl wait
  fi

  if [[ "$conjurType" == "all" || "$conjurType" == "enterprise" ]]; then
    exec_on cuke-master /opt/conjur/evoke/bin/wait_for_conjur

    exec_on cuke-master conjur authn login -u admin -p secret
    exec_on cuke-master conjur variable create existent-variable-with-undefined-value
    exec_on cuke-master conjur variable create existent-variable-with-defined-value
    exec_on cuke-master conjur variable values add existent-variable-with-defined-value existent-variable-defined-value
  fi
}

getKeys() {
  local conjurType="$1"

  if [[ "$conjurType" == "enterprise" ]]; then
    exec_on cuke-master conjur user rotate_api_key
  elif [[ "$conjurType" == "oss" ]]; then
    exec_on conjur conjurctl role retrieve-key cucumber:user:${CONJUR_AUTHN_LOGIN:-admin}
  fi
}

getCert() {
  exec_on cuke-master cat /opt/conjur/etc/ssl/ca.pem
}
