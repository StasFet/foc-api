#!/bin/bash

docker compose up --build --watch
echo "Docker container shut down!"
docker image prune -f         # removes dangling images