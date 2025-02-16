.PHONY: all windows-amd64 linux-amd64 clean init

# Binary name and version
BINARY_NAME=nav
VERSION?=1.0.0

# Output directories and binary names
WINDOWS_AMD64=$(BINARY_NAME)_windows_amd64.exe
LINUX_AMD64=$(BINARY_NAME)_linux_amd64

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -w -s"

# Build all targets
all: windows-amd64 linux-amd64

# Windows builds
windows-amd64:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/windows_amd64/$(WINDOWS_AMD64) main.go

# Linux builds
linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/linux_amd64/$(LINUX_AMD64) main.go

# Build for current platform only
build-current:
	go build $(LDFLAGS) -o build/$(BINARY_NAME) main.go

# Clean build directory
clean:
	rm -rf build/

# Create build directories
init:
	mkdir -p build/windows_amd64
	mkdir -p build/linux_amd64

# Create release archives
package: all
	cd build/windows_amd64 && zip -r ../$(BINARY_NAME)_$(VERSION)_windows_amd64.zip .
	cd build/linux_amd64 && tar czf ../$(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz .