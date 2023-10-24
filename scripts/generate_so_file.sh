#!/bin/bash

# TODO: hase some issue with tagging the image
# TODO: and take .so file coz of the container name
# TODO: make a cache for the image

docker build --no-cache -t mtso:0.0.3 . -f ./docker/deb/Dockerfile

# Start a container in the background
docker run -d --name mt_db mtso:0.0.3

# Copy the built .so file from the container to the local bin directory
docker cp mt_db:/app/bin/mt_deb.so ./bin/

# Stop and remove the container
docker rm -f mt_db