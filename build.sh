#!/bin/bash -e

function repo_root() {
  git rev-parse --show-toplevel
}

function project_version() {
  # VERSION derived from CHANGELOG and automated release library
  echo "$(<"$(repo_root)/VERSION")"
}

function project_semantic_version() {
  local version
  version="$(project_version)"

  # Remove Jenkins build number from VERSION
  echo "${version/-*/}"
}

CURRENT_DIR=$(pwd)

echo "Current dir: $CURRENT_DIR"

MOUNT_DIR="/summon-conjur"

GORELEASER_IMAGE="goreleaser/goreleaser:latest"

VERSION="$(project_semantic_version)"

docker pull "${GORELEASER_IMAGE}"
docker run --rm -t \
  --env GITHUB_TOKEN \
  --env GOTOOLCHAIN=auto \
  --env VERSION="${VERSION}" \
  --entrypoint "/sbin/tini" \
  -v "$CURRENT_DIR:$MOUNT_DIR" \
  -w "$MOUNT_DIR" \
  "${GORELEASER_IMAGE}" \
  -- sh -c "git config --global --add safe.directory $MOUNT_DIR && \
    /entrypoint.sh --clean $@ && \
    rm ./dist/goreleaser/artifacts.json"

echo "Releases built. Archives can be found in dist/goreleaser"
