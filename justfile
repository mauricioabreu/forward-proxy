# Run proxy server
run:
    go run main.go --forbidden-hosts forbidden_hosts.txt

# Execute test suite
test:
    go test -v ./...

# Check code quality
lint:
    golangci-lint run -v ./...
