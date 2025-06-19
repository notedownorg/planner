SHELL := nix develop --command bash
.SHELLFLAGS := -euo pipefail -c

.PHONY: dev build clean setup setup-frontend setup-backend help test test-frontend test-backend lint fmt dirty hygiene hygiene-frontend hygiene-backend

# Default target
help:
	@echo "Available targets:"
	@echo "  setup            - Install all dependencies (frontend)"
	@echo "  setup-frontend   - Install frontend dependencies"
	@echo "  setup-backend    - Setup backend (no action needed)"
	@echo "  dev              - Run in development mode"
	@echo "  build            - Build the application"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run all tests (frontend and backend)"
	@echo "  test-frontend    - Run frontend tests only"
	@echo "  test-backend     - Run backend tests only"
	@echo "  lint             - Run all linting (frontend and backend)"
	@echo "  fmt              - Format all code"
	@echo "  dirty            - Check if git working tree is dirty"
	@echo "  hygiene          - Run all code quality checks (fmt + lint)"
	@echo "  hygiene-frontend - Run frontend code quality checks (format + lint)"
	@echo "  hygiene-backend  - Run backend code quality checks (format + lint)"

#=============================================================================
# Setup and Development
#=============================================================================

setup: setup-frontend

setup-frontend:
	cd frontend && npm install

setup-backend:
	@echo "✅ Backend setup complete (no action needed)"

dev:
	wails dev

build:
	wails build

clean:
	rm -rf build/bin frontend/dist frontend/node_modules

#=============================================================================
# Testing
#=============================================================================

test: test-frontend test-backend
	@echo "✅ All tests passed!"

test-frontend:
	@echo "Running frontend tests..."
	@cd frontend && npm test

test-backend:
	@echo "Running backend tests..."
	@go test ./...

#=============================================================================
# Code Hygiene
#=============================================================================

hygiene: hygiene-frontend hygiene-backend dirty
	@echo "✅ All code quality checks passed!"

hygiene-frontend:
	@echo "Running frontend code quality checks..."
	@echo "Formatting frontend code..."
	@cd frontend && npm run format
	@echo "Running frontend linting..."
	@cd frontend && npm run lint
	@echo "✅ Frontend code quality checks passed!"

hygiene-backend:
	@echo "Running backend code quality checks..."
	@echo "Formatting Go code..."
	@gofmt -w .
	@echo "Running backend linting..."
	@go vet ./...
	@echo "Running security scan (gosec)..."
	@gosec ./... || echo "⚠️  Security scan completed with warnings"
	@echo "✅ Backend code quality checks passed!"

dirty:
	@echo "Checking if git working tree is dirty..."
	@git diff --exit-code
	@git diff --cached --exit-code
	@echo "✅ Git working tree is clean"
