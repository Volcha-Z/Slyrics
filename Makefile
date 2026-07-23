.PHONY: build test

build:
	go build -ldflags '-w -s' -o slyrics .

test:
	go test ./...
