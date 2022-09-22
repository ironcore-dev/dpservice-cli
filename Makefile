all: build

# Build binary
build:
	go build -o bin/dsclient main.go
