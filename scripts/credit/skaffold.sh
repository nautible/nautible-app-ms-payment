#!/bin/bash

cd ../../
skaffold dev --filename=./scripts/credit/skaffold.yaml --no-prune=false --cache-artifacts=false $1
