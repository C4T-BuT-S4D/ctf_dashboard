#!/bin/bash -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)"
COMMAND="$*"

# shellcheck disable=SC2016
VERSION='v$$NEO_VERSION$$'

cat >"${DIR}/client_config.yml" <<-EOT
host: "$$NEO_ADDR$$"
exploit_dir: "data"
grpc_auth_key: "$$PASSWORD$$"
EOT

IMAGE="ghcr.io/c4t-but-s4d/neo_env:${VERSION}"
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
    --network host \
    --name "${CONTAINER_NAME}" \
    --hostname "${CONTAINER_NAME}" \
    "${IMAGE}" \
    "${COMMAND}"
fi
