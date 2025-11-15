.PHONY: build run tidy lint test

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o plat .

run:
	go run .

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

test:
	go test ./... -v -race -cover

