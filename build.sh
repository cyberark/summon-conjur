#!/bin/bash -e

git fetch --tags  # jenkins does not do this automatically yet

docker-compose pull goreleaser

docker-compose run --rm -T --entrypoint sh goreleaser -es <<EOF
dep ensure --vendor-only
goreleaser release --rm-dist --skip-validate
EOF
