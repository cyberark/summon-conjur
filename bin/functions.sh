function startConjur() {
  local services='conjur'

  docker compose pull $services
  docker compose  up -d $services
}

exec_on() {
  local container=$1; shift
  docker exec $(docker compose ps -q $container) "$@"
}

function initEnvironment() {
  exec_on conjur conjurctl wait
}

getKeys() {
  exec_on conjur conjurctl role retrieve-key cucumber:user:${CONJUR_AUTHN_LOGIN:-admin}
}
