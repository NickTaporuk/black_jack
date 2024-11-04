# Use the official Golang image as the base image
FROM golang:1.23.2-bullseye

# Set the current working directory in the container
WORKDIR /app

# Copy everything from the current directory to the container's work directory
COPY . .

# Build the Go app
RUN go build -o bin/mt_deb.so -buildmode=plugin cmd/mtso.go && ls -lah