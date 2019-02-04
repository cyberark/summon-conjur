#!/bin/bash -e

for i in $(seq 20); do
  curl -o /dev/null -fs -X OPTIONS $CONJUR_APPLIANCE_URL > /dev/null \
    && $(curl --silent -k $CONJUR_V4_HEALTH_URL | jq .ok 2>/dev/null | grep true) \
    && break
  echo .
  sleep 2
done

# So we fail if the server isn't up yet:
curl -o /dev/null -fs -X OPTIONS $CONJUR_APPLIANCE_URL > /dev/null
$(curl --silent -k $CONJUR_V4_HEALTH_URL | jq .ok 2>/dev/null | grep true)
