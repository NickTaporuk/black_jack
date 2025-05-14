# Use Debian 11 Bullseye as the base image
FROM --platform=linux/amd64 debian:bullseye

# Update and install necessary packages for building Go code
RUN apt-get update && apt-get install -y \
    wget \
    tar \
    ca-certificates \
    g++ \
    gcc \
    libc6-dev \
    make \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Download and install Golang 1.24.3
RUN cd /tmp && wget https://golang.org/dl/go1.24.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz && \
    rm /tmp/go1.24.3.linux-amd64.tar.gz

# Set the environment variable for Go
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOBIN="$GOPATH/bin"

# Verify Go installation (This is for debugging purposes)
RUN go version

# Set environment variables for CGO (ensure that gcc and g++ are available)
ENV CC=/usr/bin/gcc
ENV CXX=/usr/bin/g++

# Set the current working directory in the container
WORKDIR /app

# Copy everything from the current directory to the container's work directory
COPY . .

# Build the Go app (Disable CGO if necessary for debugging)
RUN CGO_ENABLED=1 GOARCH=amd64 go build -v -o bin/mt_deb.so -buildmode=plugin cmd/mtso.go && ls -lah

# Optionally, disable CGO (Uncomment the line below for testing)
# RUN CGO_ENABLED=0 GOARCH=amd64 go build -v -o bin/mt_deb.so -buildmode=plugin cmd/mtso.go && ls -lah