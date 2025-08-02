#!/bin/bash
#export DOCKER_BUILDKIT=1
docker compose up --build --watch
echo "Docker container shut down!"
docker image prune -f         # removes dangling images