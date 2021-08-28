#!/bin/bash -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)"
COMMAND="$*"

# shellcheck disable=SC2016
VERSION='v$$VERSION$$'

cat >"${DIR}/client_config.yml" <<-EOT
host: "$$HOST$$"
exploit_dir: "data"
grpc_auth_key: "$$AUTH_KEY$$"
EOT

IMAGE="ghcr.io/pomo-mondreganto/neo_env:${VERSION}"
CONTAINER_NAME="neo-${VERSION}"

mkdir -p "${DIR}/exploits"
mkdir -p "${DIR}/data"

echo "Using image: ${IMAGE}"

OUT=$(docker ps --filter "name=${CONTAINER_NAME}" --format "{{ .Names }}")

if [[ $OUT ]]; then
  echo "Container already exists"
  # shellcheck disable=SC2068
  docker exec -it "${CONTAINER_NAME}" ${COMMAND[@]}
else
  echo "Starting a new container"
  docker run -it \
    --rm \
    --volume "${DIR}":/work \
    --security-opt seccomp=unconfined \
    --security-opt apparmor=unconfined \
    --cap-add=NET_ADMIN \
    --privileged \
    --name "${CONTAINER_NAME}" \
    --hostname "${CONTAINER_NAME}" \
    "${IMAGE}" \
    "${COMMAND}"
fi
