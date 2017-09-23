#!/bin/sh
go build -ldflags="-s -w" -o large main.go
rm -f usock2wsock
upx --ultra-brute -o usock2wsock large && rm -f large
