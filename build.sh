#!/bin/sh
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o usock2wsock_amd64 main.go
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o usock2wsock_arm64 main.go
lipo -create -output usock2wsock usock2wsock_amd64 usock2wsock_arm64
