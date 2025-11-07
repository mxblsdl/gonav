package main

import (
	"embed"
	"fmt"
	"os"
	"path"
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
		subEntries, _ := embeddedFiles.ReadDir("binary/" + entry.Name())
		fmt.Printf("- %s/\n", entry.Name())
		for _, subEntry := range subEntries {
			fmt.Printf("  └── %s\n", subEntry.Name())
		}
	}
}

func getInstallConfig() (destDir, binName string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("unable to find home directory: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "AppData", "Local", "nav"), "nav.exe", nil
	case "linux", "darwin":
		return filepath.Join(homeDir, ".local", "bin"), "nav", nil
	default:
		return "", "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func main() {
	destDir, binName, err := getInstallConfig()
	if err != nil {
		fmt.Printf("Error getting install config: %v\n", err)
		os.Exit(1)
	}

	// Create destination directory if it doesn't exist (Windows)
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create destination directory: %v\n", err)
		os.Exit(1)
	}

	// Read the embedded binary for current OS
	osPath := path.Join("binary", runtime.GOOS, binName)
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

	fmt.Printf("Successfully installed %s to %s\n", binName, destPath)
}
