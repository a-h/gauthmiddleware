init:
	go get ./...

test:
	go test ./...

build:
	go build ./...

.PHONY: init, test, build
