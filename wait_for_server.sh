#!/bin/bash -e

for i in $(seq 20); do
  curl -o /dev/null -fs -X OPTIONS $CONJUR_APPLIANCE_URL > /dev/null && break
  echo .
  sleep 2
done

# So we fail if the server isn't up yet:
curl -o /dev/null -fs -X OPTIONS $CONJUR_APPLIANCE_URL > /dev/null
