package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed all:binary/*
//go:embed all:binary/windows/*
//go:embed all:binary/linux/*
//go:embed all:binary/darwin/*

var embeddedFiles embed.FS

func init() {
	// Debug: List all embedded files
	entries, err := embeddedFiles.ReadDir("binary")
	if err != nil {
		fmt.Printf("Error reading binary directory: %v\n", err)
		return
	}

	fmt.Println("Embedded files in binary/:")
	for _, entry := range entries {
		subEntries, _ := embeddedFiles.ReadDir(filepath.Join("binary", entry.Name()))
		fmt.Printf("- %s/\n", entry.Name())
		for _, subEntry := range subEntries {
			fmt.Printf("  └── %s\n", subEntry.Name())
		}
	}
}

var (
	binaryName = "nav"
)

func main() {
	var destDir string
	var binName string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Unable to find home directory: %v\n", err)
	}

	switch runtime.GOOS {
	case "windows":
		destDir = filepath.Join(homeDir, "AppData", "Local", "nav")
		binName = binaryName + ".exe"
	case "linux":
		destDir = filepath.Join(homeDir, ".local", "bin")
		binName = binaryName
	case "darwin":
		destDir = filepath.Join(homeDir, ".local", "bin")
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
		fmt.Printf("Attemped to read from: %v\n", osPath)
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

// func ExpandPath(path string) string {
// 	if path[:2] == "~/" {
// 		homedir, err := os.UserHomeDir()
// 		if err != nil {
// 			fmt.Println("Error getting home directory:", err)
// 			return path
// 		}
// 		return homedir + path[1:]
// 	}
// 	return path
// }
