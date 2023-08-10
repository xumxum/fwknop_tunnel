#!/bin/bash
#go build -ldflags="-X 'test1/version.Version=1.0.0'"

# STEP 1: Determinate the required values
PACKAGE="fwknop_tunnel"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y/%m/%d %H:%M:%S')

# STEP 2: Build the ldflags

LDFLAGS=(
  "-X '${PACKAGE}/version.CommitHash=${COMMIT_HASH}'"
  "-X '${PACKAGE}/version.BuildTimestamp=${BUILD_TIMESTAMP}'"
)

# STEP 3: Actual Go build process

go build -ldflags="${LDFLAGS[*]}"

#./fwknop_tunnel --version

#mv ./fwknop_tunnel /usr/local/bin/
#systemctl restart fwknop_tunnel.service

