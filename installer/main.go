package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed binary/windows/*
//go:embed binary/linux/*
var embeddedFiles embed.FS

var (
	binaryName = "nav"
)

func main() {
	var destDir string
	var binName string

	switch runtime.GOOS {
	case "windows":
		destDir = "C:\\Program Files\\nav"
		binName = binaryName + ".exe"
	case "linux":
		destDir = ExpandPath("~/bin")
		binName = binaryName	
	}
	// Create destination directory if it doesn't exist (Windows)
	if runtime.GOOS == "windows" {
		err := os.MkdirAll(destDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create destination directory: %v\n", err)
			os.Exit(1)
		}
	}
	// Read the embedded binary for current OS
	osPath := filepath.Join("binary", runtime.GOOS, binName)
	binaryData, err := embeddedFiles.ReadFile(osPath)
	if err != nil {
		fmt.Printf("Error reading embedded binary: %v\n", err)
		os.Exit(1)
	}

	// Write the binary to destination
	destPath := filepath.Join(destDir, binName)
	err = os.WriteFile(destPath, binaryData, 0755)
	if err != nil {
		fmt.Printf("Failed to write binary to %s: %v\n", destPath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully installed %s to %s\n", binaryName, destPath)
}

func ExpandPath(path string) string {
	if path[:2] == "~/" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return path
		}
		return homedir + path[1:]
	}
	return path
}