#!/bin/bash
BIN=server

CGO_ENABLED=0 go build -o $BIN -a -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"'

docker build . -t $1
