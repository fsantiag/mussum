#!/usr/bin/env bash
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
docker push fsantiag/mussum:1:0:0
