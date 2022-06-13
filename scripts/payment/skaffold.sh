#!/bin/bash

# $1 = aws|azure

cd ../../
skaffold dev -p $1 --filename=./scripts/payment/skaffold.yaml --no-prune=false --cache-artifacts=false
