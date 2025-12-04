# Detect OS and Architecture if not provided
OS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH ?= $(shell uname -m)

# Default binary names
BINARY_LINUX = mersenne_twister_linux
BINARY_MACOS = mersenne_twister_darwin

# Compiler flags
CFLAGS = -fPIC

# Set file extensions based on OS
ifeq ($(OS), darwin)
    LIBRARY = libmtgenerator_macos.dylib
    LDFLAGS = -dynamiclib
    CGO_FLAGS = DYLD_LIBRARY_PATH=.
    BINARY = $(BINARY_MACOS)
else
    LIBRARY = libmtgenerator_linux.so
    LDFLAGS = -shared
    CGO_FLAGS = LD_LIBRARY_PATH=.
    BINARY = $(BINARY_LINUX)
endif

# Default target
all: build

# Build the shared library and Go binary
build:
	g++ $(LDFLAGS) -o $(LIBRARY) mt_generator.cpp $(CFLAGS)
	#g++ -dynamiclib -o libmtgenerator_macos.dylib mt_generator.cpp -fPIC
	CGO_ENABLED=1 GOOS=$(OS) GOARCH=$(ARCH) $(CGO_FLAGS) go build -o $(BINARY) mt_wrapper_darwin.go

# Run the application
run: build
	$(CGO_FLAGS) ./$(BINARY)

# Cross-compilation for Linux and macOS
build-linux:
	g++ -shared -o libmtgenerator_linux.so mt_generator.cpp -fPIC
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 LD_LIBRARY_PATH=. go build -o $(BINARY_LINUX) mt_wrapper_linux.go

build-macos:
	g++ -dynamiclib -o libmtgenerator_macos.dylib mt_generator.cpp -fPIC
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 DYLD_LIBRARY_PATH=. go build -o $(BINARY_MACOS) mt_wrapper_darwin.go

# Clean build artifacts
clean:
	rm -f $(BINARY_LINUX) $(BINARY_MACOS) libmtgenerator_linux.so libmtgenerator_macos.dylib
