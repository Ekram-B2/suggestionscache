#!/bin/bash

set -e -x

MYWORKDIR=$PWD
DOCKER_GOPATH=/home/alpine/go
LOCAL_GOPATH=/var/go
LOCAL_VOLUME=/Users/ekrambhuiyan/go


# build the godeps.txt
docker run -i --rm -v "${LOCAL_VOLUME}:${DOCKER_GOPATH}" -e "GOPATH=${DOCKER_GOPATH}" -w "${DOCKER_GOPATH}/src/github.com/ezoic/gvlcache" superbaddude/golang:1.12.13-alpine \
   /bin/bash <<'EOF'
./build.sh
EOF
