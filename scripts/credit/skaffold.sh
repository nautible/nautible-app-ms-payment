#!/bin/bash

cd ../../
skaffold dev --filename=./scripts/credit/skaffold.yaml --port-forward
