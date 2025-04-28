# Define the binary name
BINARY_NAME=redez

# Define the source directory
SRC_DIR := src

# Default target
all: build

# Build the Go application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./src #Build the Go application
	@echo "Build complete: $(BINARY_NAME)"
	@echo ""

#Build the Redez-cli
build-cli:
	@echo "Building $(BINARY_NAME)-cli..."
	@go build -o $(BINARY_NAME)-cli ./cli/main.go
	@echo "Build complete: $(BINARY_NAME)-cli"
	@echo ""

# Run the application
run: build
	@echo "Starting $(BINARY_NAME)..."
	@./$(BINARY_NAME)
	@echo ""

# Run the application-cli
run-cli: build build-cli
	@echo "Starting $(BINARY_NAME)-cli..."
	@./$(BINARY_NAME)-cli
	@echo ""

# Clean the build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-cli $(BINARY_NAME).aof
	@echo "Clean complete."

.PHONY: all build run clean
