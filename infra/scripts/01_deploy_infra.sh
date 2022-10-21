#!/bin/bash

###################
### Infra Setup ###
###################

kind create cluster \
  --name test \
  --config kind-config.yaml \
  --image=kindest/node:v1.24.0
