#!/bin/bash

docker build --build-arg HTTP_PROXY=$http_proxy --build-arg HTTPS_PROXY=$http_proxy -t $IMAGE -f ./package/cash/Dockerfile .
