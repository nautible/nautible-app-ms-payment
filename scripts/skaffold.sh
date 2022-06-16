#!/bin/bash

# $1 = aws|azure

cd ../
skaffold dev -p $1 --filename=./scripts/skaffold.yaml
