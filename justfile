# Run proxy server
run:
    go run main.go

# Execute test suite
test:
    go test -v ./...

# Check code quality
lint:
    golangci-lint run -v ./...
