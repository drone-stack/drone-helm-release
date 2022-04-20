#!/bin/bash

# disable cgo
export CGO_ENABLED=0

set -e
set -x

go build -o release/linux/amd64/plugin
GOOS=linux GOARCH=amd64 go build -o release/linux/amd64/plugin-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o release/darwin/amd64/plugin-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o release/darwin/arm64/plugin-darwin-arm64
