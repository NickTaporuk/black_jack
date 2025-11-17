LUA_VERSION=5.4.7
LUA_PATH=/opt/homebrew/Cellar/lua/$(LUA_VERSION)

CC=gcc
CFLAGS=-I$(LUA_PATH)/include/lua
LDFLAGS=-L$(LUA_PATH)/lib -llua



##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ##  Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build
.PHONY: build
build: ##  Build the project locally
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64  go build -o bin/mt_darwin.so -buildmode=plugin ./cmd/mtso.go

.PHONY: build-debug-mode
build-debug-mode: ## This build will not be optimized to run just only debug mode
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64  go build  -gcflags="all=-N -l" -o bin/mt_darwin_debug.so -buildmode=plugin ./cmd/mtso.go

.PHONY: build-debian-debug-mode
build-debian-debug-mode: ## This build will not be optimized to run just only debug mode
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build  -gcflags="all=-N -l" -o bin/mt_linux_debug.so -buildmode=plugin ./cmd/mtso.go


.PHONY: build-linux
build-linux: ##  Build the project locally
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64  go build -o bin/mt_deb.so -buildmode=plugin ./cmd/mtso.go

build-docker-deb:  ## Build the project .so file for debian linux by docker
	./scripts/generate_so_file.sh

build-lua-so-go:
CGO_ENABLED=1 go build -o mtso_lua \
    -gcflags "all=-N -l" \
    -tags lua \
    -ldflags "`pkg-config --libs lua5.4`" \
    -cgo CFLAGS="`pkg-config --cflags lua5.4`" cmd/mtso_lua.go


all: build-lua-so

build-lua-so:
	$(CC) -shared -o libjackpot.so jackpot.c $(CFLAGS) $(LDFLAGS) -fPIC

.PHONY: build-debian-c++_so
build-debian-c++_so: ## This is the building .so file from c++ to golang linux part
	cd mtwrapper && g++ -std=c++17 -shared -fPIC -O3 \
		-o libmtgenerator_linux.so mt_generator.cpp
