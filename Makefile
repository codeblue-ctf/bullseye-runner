GOOS := linux
GOARCH := amd64

BUILD := $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X=main.Build=$(BUILD)"

.PHONY: build
build: build-proto
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/runner $(LDFLAGS) -v ./src

.PHONY: build-proto
build-proto: ./proto/*.proto
	protoc -I ./proto --go_out=./src/proto ./proto/*.proto

.PHONY: clean
clean:
	rm -f bin/*