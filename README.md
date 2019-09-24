# bulls-eye-runner
New Bull's Eye Runner implemented in Golang

## Overview
- runner-master
  - scheduler for Bull's Eye
  - send gRPC request to runner-worker to run evaluation
  - be cancellable with gRPC connection
- runner-worker
  - evaluate docker-compose.yml
  - scalabe

## Prerequisites
- Golang (>= 1.12)
- Make
- protobuf

## Usage
run `make`

## TODO
- scheduler in runner-master
- forwarding X11 display to show audience what's going on
- send evaluation result in real-time to bullseye-web