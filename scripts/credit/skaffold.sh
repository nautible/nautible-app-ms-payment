#!/bin/bash

cd ../../
skaffold dev --filename=./scripts/credit/skaffold.yaml $1
