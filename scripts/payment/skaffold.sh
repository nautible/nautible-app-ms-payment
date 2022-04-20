#!/bin/bash

cd ../../
skaffold dev --filename=./scripts/payment/skaffold.yaml --no-prune=false --cache-artifacts=false $1
