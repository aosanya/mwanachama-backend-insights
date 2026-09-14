.PHONY: build test vet clean

## Verify the module compiles cleanly.
build:
	go build ./...

## Unit tests (sqlite-backed, no DB required).
test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -rf bin/
