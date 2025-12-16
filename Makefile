.PHONY: build build-web build-server clean

# Variables
BINARY_NAME=yan
OUTPUT_DIR=bin
GO_MAIN=./cmd/root
WEB_DIR=web

# Build web frontend
build-web:
	@echo "Building web frontend..."
	cd $(WEB_DIR) && pnpm install && pnpm build-only
	@echo "Web frontend build completed"

# Build Go server with static linking
build-server:
	@echo "Building Go server..."
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o $(OUTPUT_DIR)/$(BINARY_NAME) $(GO_MAIN)
	@echo "Go server build completed: $(OUTPUT_DIR)/$(BINARY_NAME)"

# Full build (web first, then server)
build: build-web build-server
	@echo "All build tasks completed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(OUTPUT_DIR)
	cd $(WEB_DIR) && rm -rf dist
	@echo "Clean completed"
