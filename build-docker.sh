#! /bin/bash

# build collector image
docker build -t adrisongomez/custom-collector . -f ./cmd/collectors/Dockerfile 

# build server image
docker build -t adrisongomez/server . -f ./cmd/server/Dockerfile 

docker push adrisongomez/custom-collector
docker push adrisongomez/server