#! /bin/bash

# build collector image
docker build -t custom-collector . -f ./cmd/collectors/Dockerfile 

# build server image
docker build -t server . -f ./cmd/server/Dockerfile 