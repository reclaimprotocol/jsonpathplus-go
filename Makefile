# Makefile for JSONPath Plus Go

# Variables
BINARY_NAME=jsonpath-example
PACKAGE_NAME=jsonpathplus-go
GO_VERSION=$(shell go version | cut -d' ' -f3)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
.PHONY: help
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the example binary
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p bin
	go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/examples/basic
	@echo "Binary built: bin/${BINARY_NAME}"

## build-all: Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 ./cmd/examples/basic
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 ./cmd/examples/basic
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 ./cmd/examples/basic
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe ./cmd/examples/basic
	@echo "All binaries built in bin/"

## test: Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Coverage report: coverage.out"

## test-coverage: Run tests with coverage report
.PHONY: test-coverage
test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## bench: Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

## lint: Run linters
.PHONY: lint
lint:
	@echo "Running linters..."
	golangci-lint run ./...

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	gofmt -s -w .

## vet: Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

## mod-tidy: Tidy go modules
.PHONY: mod-tidy
mod-tidy:
	@echo "Tidying go modules..."
	go mod tidy
	go mod verify

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache -testcache -modcache

## security: Run security checks
.PHONY: security
security:
	@echo "Running security checks..."
	gosec ./...

## deps: Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download

## update-deps: Update all dependencies
.PHONY: update-deps
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

## install: Install the package
.PHONY: install
install:
	@echo "Installing package..."
	go install ${LDFLAGS} ./cmd/examples/basic

## run: Run the example
.PHONY: run
run:
	@echo "Running example..."
	go run ./cmd/examples/basic 2>/dev/null

## run-production: Run the production example
.PHONY: run-production
run-production:
	@echo "Running production example..."
	go run ./cmd/examples/production

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t ${PACKAGE_NAME}:latest .

## docker-run: Run Docker container
.PHONY: docker-run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 ${PACKAGE_NAME}:latest

## check: Run all checks (fmt, vet, lint, test)
.PHONY: check
check: fmt vet lint test
	@echo "All checks passed!"

## ci: Run CI pipeline locally
.PHONY: ci
ci: deps mod-tidy check security bench
	@echo "CI pipeline completed successfully!"

## release: Prepare for release
.PHONY: release
release: clean ci build-all
	@echo "Release artifacts prepared in bin/"

## stats: Show project statistics
.PHONY: stats
stats:
	@echo "Project Statistics:"
	@echo "=================="
	@echo "Go version: ${GO_VERSION}"
	@echo "Git commit: ${GIT_COMMIT}"
	@echo "Build time: ${BUILD_TIME}"
	@echo ""
	@echo "Code statistics:"
	@find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1
	@echo ""
	@echo "Test coverage:"
	@go test -coverprofile=/tmp/coverage.out ./... >/dev/null 2>&1 && go tool cover -func=/tmp/coverage.out | tail -1 || echo "No coverage data"

## serve-docs: Serve documentation locally
.PHONY: serve-docs
serve-docs:
	@echo "Serving documentation on http://localhost:6060"
	godoc -http=:6060 -play

## profile: Run CPU profiling
.PHONY: profile
profile:
	@echo "Running CPU profiling..."
	go run -cpuprofile=cpu.prof ./cmd/examples/production
	@echo "Profile saved to cpu.prof"

## memory-profile: Run memory profiling
.PHONY: memory-profile
memory-profile:
	@echo "Running memory profiling..."
	go run -memprofile=mem.prof ./cmd/examples/production
	@echo "Profile saved to mem.prof"