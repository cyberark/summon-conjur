#!/bin/bash -e

CURRENT_DIR=$(pwd)
MOUNT_DIR="/summon-conjur"

GORELEASER_IMAGE="goreleaser/goreleaser:latest"
GORELEASER_ARGS="--rm-dist --snapshot"

git fetch --tags  # jenkins does not do this automatically yet

docker pull "${GORELEASER_IMAGE}"
docker run --rm -t \
  -v "$CURRENT_DIR:$MOUNT_DIR" \
  -w "$MOUNT_DIR" \
  "${GORELEASER_IMAGE}" ${GORELEASER_ARGS}
