.PHONY: test cover build install

# Run all tests
test:
	go test ./...

# Run tests with coverage report
cover:
	go test -cover ./...

# Run tests with coverage profile (for coverage percentage and HTML report)
cover-profile:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Optional: HTML coverage report (open in browser)
cover-html: cover-profile
	go tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in your browser"

# Build the CLI binary
build:
	go build -o go-gen-r ./cmd/go-gen-r

# Install the CLI to $GOPATH/bin or $GOBIN
install: build
	cp go-gen-r /usr/local/bin/go-gen-r 2>/dev/null || true
