#!/bin/bash

mkdir -p installer/binary
mkdir -p installer/dist

echo "Building Linux Binary"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o installer/binary/linux/nav

echo "Building Darwin binary"
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o installer/binary/darwin/nav
chmod +x installer/binary/darwin/* installer/binary/linux/*

echo "Building Windows exe"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o installer/binary/windows/nav.exe

echo "Building installers"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o installer/dist/nav-linux-amd64 installer/main.go
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o installer/dist/nav-darwin-amd64 installer/main.go
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o installer/dist/nav-windows-amd64.exe installer/main.go

chmod +x installer/dist/nav-linux** installer/dist/nav-darwin**

echo "Build complete. Binaries located in installer/binary/"

