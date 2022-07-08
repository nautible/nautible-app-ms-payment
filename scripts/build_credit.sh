#!/bin/bash

#aws|azure
CLOUD=$1

docker build --build-arg CLOUD=$CLOUD --build-arg HTTP_PROXY=$http_proxy --build-arg HTTPS_PROXY=$http_proxy -t $IMAGE -f ./package/credit/Dockerfile.local .
