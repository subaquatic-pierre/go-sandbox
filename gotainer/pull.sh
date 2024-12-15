#!/bin/bash

set -e


# Use default image or first argument passed to script
container=$(docker create "$image")

mkdir -p "./assets/${image}"


# Create tar of exported container
docker export "$container" -o "./assets/${image}/${image}.tar.gz" > /dev/null

# Remove container after creating tar
docker rm "$container" > /dev/null

# Write name of command to direcory
docker inspect -f '{{.Config.Cmd}}' "$image:latest" | tr -d '[]\n' > "./assets/${image}/${image}-cmd"
