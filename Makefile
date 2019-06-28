GOOS := linux
GOARCH := amd64

BUILD := $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X=main.Build=$(BUILD)"

.PHONY: build
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/runner $(LDFLAGS) -v ./src

.PHONY: clean
clean:
	rm -f bin/*