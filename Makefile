all: build

build:
	@echo "Building..."
	@go build -o verctl ./cmd/verctl

build-snapshot:
	@echo "Building snapshot..."
	@goreleaser build --single-target --snapshot --clean

release:
	@goreleaser release --clean

test:
	@echo "Testing..."
	@go test ./... -v

coverage:
	@echo "Running coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out

fmt:
	@echo "Formatting..."
	@go fmt ./...

lint:
	@echo "Linting..."
	@gofmt -l .
	@go vet ./...

clean:
	@echo "Cleaning..."
	@rm -rf verctl dist/ coverage.out coverage.html

.PHONY: all build build-snapshot release test coverage fmt lint clean
