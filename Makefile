.PHONY: build test lint lint-example clean

# Build standalone CLI tool
build:
	go build -o loglinter ./cmd/loglinter

# Build custom golangci-lint binary with plugin included
# Requires existing golangci-lint installation
build-custom-gcl:
	golangci-lint custom

# Run all tests
test:
	go test -v ./...

# Run linter on example file using standalone binary
lint-example: build
	./loglinter ./testdata/src/example || true

# Run golangci-lint on the project itself
lint:
	golangci-lint run

clean:
	rm -f loglinter custom-gcl
