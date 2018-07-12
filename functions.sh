startConjur() {
  # Start a local Conjur to test against
  docker-compose pull postgres conjur cuke-master
  docker-compose up -d conjur cuke-master
}

exec_on() {
  local container=$1; shift
  docker exec $(docker-compose ps -q $container) "$@"
}


initEnvironment() {
  # Delay to allow time for Conjur to come up
  # TODO: remove this once we have HEALTHCHECK in place
  exec_on conjur conjurctl wait
  exec_on cuke-master /opt/conjur/evoke/bin/wait_for_conjur

  exec_on cuke-master conjur authn login -u admin -p secret
  exec_on cuke-master conjur variable create existent-variable-with-undefined-value
  exec_on cuke-master conjur variable create existent-variable-with-defined-value
  exec_on cuke-master conjur variable values add existent-variable-with-defined-value existent-variable-defined-value
}

getKeys() {
  exec_on conjur conjurctl role retrieve-key cucumber:user:${CONJUR_AUTHN_LOGIN:-admin}
  exec_on cuke-master conjur user rotate_api_key
}

getCert() {
  exec_on cuke-master cat /opt/conjur/etc/ssl/ca.pem
}
