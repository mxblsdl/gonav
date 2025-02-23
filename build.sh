#!/bin/bash

mkdir -p installer/binary

echo "Building Linux Binary"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o installer/binary/linux/nav

echo "Building Darwin binary"
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o installer/binary/darwin/nav

echo "Building Windows exe"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o installer/binary/windows/nav.exe
# # Set executable permissions for Unix-based systems
# Set executable permissions
chmod +x installer/binary/linux/*

echo "Build complete. Binaries located in installer/binary/"
ls -l -R installer/binary/