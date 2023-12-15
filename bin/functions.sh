function startConjur() {
  local services='conjur'

  docker compose $COMPOSE_ARGS pull $services
  docker compose $COMPOSE_ARGS up -d $services
}

exec_on() {
  local container=$1; shift
  docker exec $(docker compose $COMPOSE_ARGS ps -q $container) "$@"
}

function initEnvironment() {
  exec_on conjur conjurctl wait
}

getKeys() {
  exec_on conjur conjurctl role retrieve-key cucumber:user:${CONJUR_AUTHN_LOGIN:-admin}
}
