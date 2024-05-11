#!/usr/bin/env bash

echo "Start packaging..."

go mod download
go mod tidy

rm -rf ./out
mkdir ./out

generate(){
    CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -o harbor_$1_$2 -ldflags '-s -w'
    mv harbor_$1_$2 ./out/
}

echo "[1] Linux from amd64"
generate linux amd64
echo "[2] Linux from arm64"
generate linux arm64