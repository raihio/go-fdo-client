# Variables
BINARY_NAME = fdo_client
BUILD_DIR = ./cmd

# Build the Go project
build:
	go build -o $(BINARY_NAME) $(BUILD_DIR)

# Clean up the binary
clean:
	rm -f $(BINARY_NAME)

# Default target
all: build
