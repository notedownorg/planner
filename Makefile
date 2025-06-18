.PHONY: dev build clean setup help

# Default target
help:
	@echo "Available targets:"
	@echo "  setup  - Install npm dependencies"
	@echo "  dev    - Run in development mode"
	@echo "  build  - Build the application"
	@echo "  clean  - Clean build artifacts"

# Setup dependencies
setup:
	cd frontend && nix develop -c npm install

# Development mode
dev:
	nix develop -c wails dev

# Build the application
build:
	nix develop -c wails build

# Clean build artifacts
clean:
	rm -rf build/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules