#!/usr/bin/env sh

set -e

if [ -z ${IMAGE_VERSION} ]; then
    echo "IMAGE_VERSION env var needs to be set"
    exit 1
fi

DIR="$( cd "$( dirname "${0}" )" && pwd )"
ROOT_DIR=${DIR}/../..
REPOSITORY="quay.io/slok"
IMAGE="tracing-example"
TARGET_IMAGE=${REPOSITORY}/${IMAGE}


docker build \
    -t ${TARGET_IMAGE}:${IMAGE_VERSION} \
    -f ${ROOT_DIR}/docker/prod/Dockerfile .
