#!/usr/bin/env bash

if [ $# != 1 ]; then
    echo "Usage: $0 [Code File Name]"
    exit 0
fi

_BIN_NAME=asg

rm ./pkg/*
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./pkg/${_BIN_NAME}_linux_amd64
upx ./pkg/${_BIN_NAME}_linux_amd64
GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o ./pkg/${_BIN_NAME}_windows_amd64.exe
upx ./pkg/${_BIN_NAME}_windows_amd64.exe
GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o ./pkg/${_BIN_NAME}_darwin_amd64
upx ./pkg/${_BIN_NAME}_darwin_amd64
