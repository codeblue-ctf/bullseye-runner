GOOS := linux
GOARCH := amd64

BUILD := $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X=main.Build=$(BUILD)"

.PHONY: build
build: build-master build-worker build-notification build-client

.PHONY: build-master
build-master: build-proto
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/runner-master $(LDFLAGS) -v ./cmd/master

.PHONY: build-worker
build-worker: build-proto
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/runner-worker $(LDFLAGS) -v ./cmd/worker

.PHONY: build-notification
build-notification:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/notification $(LDFLAGS) -v ./cmd/notification

.PHONY: build-client
build-client: build-proto
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/client $(LDFLAGS) -v ./cmd/client

.PHONY: build-proto
build-proto: ./proto/*.proto
	protoc -I ./proto --go_out=plugins=grpc:./proto ./proto/*.proto

.PHONY: clean
clean:
	rm -f bin/*
