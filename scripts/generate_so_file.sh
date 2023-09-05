docker build -t mtso . -f ./docker/Dockerfile

# Start a container in the background
docker run -d --name mtso mtso

# Copy the built .so file from the container to the local bin directory
docker cp mtso:/app/bin/mt_deb.so ./bin/

# Stop and remove the container
docker rm -f mtso