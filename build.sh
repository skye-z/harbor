#!/usr/bin/env bash

echo "Start packaging..."

go mod download
go mod tidy

rm -rf ./out
mkdir ./out

generate(){
    CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -o harbor-$1-$2 -ldflags '-s -w'
    mv harbor-$1-$2 ./out/
}

echo "[1] MacOS from amd64"
generate darwin amd64
echo "[2] MacOS from arm64"
generate darwin arm64
echo "[3] Linux from amd64"
generate linux amd64
echo "[4] Linux from arm64"
generate linux arm64