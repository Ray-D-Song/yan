.PHONY: build build-web build-server clean dev

# Variables
BINARY_NAME=yan
OUTPUT_DIR=bin
GO_MAIN=./cmd/root
WEB_DIR=web

# Detect OS and set appropriate flags
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	# Linux: Enable CGO, optionally use static linking if sqlite3 static lib is available
	CGO := 1
	LDFLAGS := -s -w
else ifeq ($(UNAME_S),Darwin)
	# macOS: Enable CGO, dynamic linking (macOS doesn't support static linking)
	CGO := 1
	LDFLAGS := -s -w
else
	# Other OS: Enable CGO with default flags
	CGO := 1
	LDFLAGS := -s -w
endif

# Build web frontend
build-web:
	@echo "Building web frontend..."
	cd $(WEB_DIR) && pnpm install && pnpm build-only
	@echo "Web frontend build completed"

# Build Go server with CGO enabled for sqlite3
build-server:
	@echo "Building Go server for $(UNAME_S) with CGO_ENABLED=$(CGO)..."
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=$(CGO) go build -trimpath -ldflags '$(LDFLAGS)' -o $(OUTPUT_DIR)/$(BINARY_NAME) $(GO_MAIN)
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

# Development mode (run both web and server in dev mode)
dev:
	@echo "Starting development servers..."
	@cd $(WEB_DIR) && pnpm dev & \
	CGO_ENABLED=$(CGO) go run $(GO_MAIN) server & \
	wait
