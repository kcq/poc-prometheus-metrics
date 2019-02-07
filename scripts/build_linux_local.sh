#!/usr/bin/env bash

set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
BDIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

echo "Preparing..."

mkdir -p ${BDIR}/bin/

echo "Building PoC..."

pushd ${BDIR}/cmd/prom-metrics-service
pwd
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags 'netgo osusergo' -o "${BDIR}/bin/prom-metrics-service"
popd
