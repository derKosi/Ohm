.PHONY: build run test lint clean install

BINARY=ohm
VERSION?=0.1.0
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/ohm/

run: build
	./$(BINARY) scan

test:
	go test ./... -v

lint:
	golangci-lint run ./...

install: build
	cp $(BINARY) /usr/local/bin/

clean:
	rm -f $(BINARY)
	rm -f ohm-cleanup-*.sh ohm-cleanup-*.ps1 ohm-cleanup-*.bat

# Cross-compile all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/ohm-linux-amd64 ./cmd/ohm/
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/ohm-linux-arm64 ./cmd/ohm/
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/ohm-darwin-amd64 ./cmd/ohm/
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/ohm-darwin-arm64 ./cmd/ohm/
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/ohm-windows-amd64.exe ./cmd/ohm/
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o dist/ohm-windows-arm64.exe ./cmd/ohm/

# Quick scan (no TUI, just output)
scan: build
	./$(BINARY) scan

# Generate script only
generate: build
	./$(BINARY) generate
