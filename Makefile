.PHONY: all verify build clean

all: clean verify build

verify:
	go fmt ./...
	go vet ./...
	go test ./...

build:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/plexname ./cmd/plexname

clean:
	rm -rf dist
