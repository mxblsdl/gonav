#!/bin/bash

set -e

# Change to project root
cd "$(dirname "$0")/.."

# Clean previous builds
rm -rf installer/binary installer/dist

# Create directory structure
mkdir -p installer/binary/{linux,darwin,windows}
mkdir -p installer/dist

echo "Building platform binaries..."
# Build the nav binaries first
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o installer/binary/linux/nav main.go
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o installer/binary/darwin/nav main.go
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o installer/binary/windows/nav.exe main.go

# Set permissions
chmod +x installer/binary/linux/nav installer/binary/darwin/nav

# Verify binaries exist and show sizes
echo "Verifying binaries..."
for bin in installer/binary/{linux/nav,darwin/nav,windows/nav.exe}; do
    if [ ! -f "$bin" ]; then
        echo "Error: Binary not created: $bin"
        ls -l installer/binary/
        exit 1
    fi
    size=$(ls -lh "$bin" | awk '{print $5}')
    echo "✓ Found $bin (size: $size)"
done

# Now build the installer after the binaries are in place
echo "Building installers..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o installer/dist/nav-windows-amd64.exe ./installer
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o installer/dist/nav-linux-amd64 ./installer
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o installer/dist/nav-darwin-amd64 ./installer

chmod +x installer/dist/nav-linux-amd64 installer/dist/nav-darwin-amd64

echo "Build complete! Installers created in installer/dist/"
ls -lh installer/dist/