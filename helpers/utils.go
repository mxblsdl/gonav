package helpers

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

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

func OpenInEditor(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Use notepad on Windows
		cmd = exec.Command("notepad", filePath)
	case "darwin":
		// Use nano on macOS
		cmd = exec.Command("nano", filePath)
	default: // Linux
		// Try common editors in order of preference
		editor := os.Getenv("EDITOR")
		if editor == "" {
			// Fallback to common editors
			for _, e := range []string{"nano", "vi", "vim"} {
				if _, err := exec.LookPath(e); err == nil {
					editor = e
					break
				}
			}
		}
		if editor == "" {
			return fmt.Errorf("no editor found")
		}
		cmd = exec.Command(editor, filePath)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error opening config file in editor:", err)
		os.Exit(1)
	}
	return nil
}

func PrintConfigMessage(hour int64, cacheFile string) {
	configFile := viper.ConfigFileUsed()
	printInterval := time.Duration(hour) * time.Hour

	lastPrinted := time.Time{}
	if data, err := os.ReadFile(cacheFile); err == nil {
		if t, err := time.Parse(time.RFC3339, string(data)); err == nil {
			lastPrinted = t
		}
	}

	if time.Since(lastPrinted) > printInterval {
		fmt.Println("Config file found:", configFile, "\nThis message is printed once every 8 hours")
		err := os.WriteFile(cacheFile, []byte(time.Now().Format(time.RFC3339)), 0644)
		if err != nil {
			fmt.Println("Error writing to cache file:", err)
		}
	}
}

func createConfig(defaultConfigPath string) {
	fmt.Printf("Config file does not exist. Do you want to create a default config file? %s(y/n)%s\n", ColorGreen, ColorReset)
	var response string
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		fmt.Println("Creating default config file at " + defaultConfigPath)

		configYaml := navConfig{
			Folders: []string{
				"~/Documents",
				"~/Projects",
			},
			MaxDepth: 3,
			Comments: "Add folders to search through in the folders section. This line can be deleted.",
		}

		data, err := yaml.Marshal(&configYaml)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = os.WriteFile(defaultConfigPath, data, 0644)
		if err != nil {
			fmt.Println("Error creating config file:", err)
			os.Exit(1)
		}
		fmt.Println("Default config file created at", defaultConfigPath)
		err = OpenInEditor(defaultConfigPath)
		if err != nil {
			fmt.Printf("%sError opening config file: %v%s\n", ColorRed, err, ColorReset)
			return
		}
	}
}

func InitConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".gonav")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		cacheFile := filepath.Join(os.TempDir(), "gonav_last_printed")
		PrintConfigMessage(8, cacheFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defaultConfigPath := filepath.Join(home, ".gonav.yaml")
		createConfig(defaultConfigPath)
	}
}

func ScanWithWalkDir(root string) ([]string, error) {
	var folders []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if strings.HasPrefix(filepath.Base(path), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Only add directories to the slice
		if d.IsDir() {
			folders = append(folders, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return folders, nil
}

func OpenShellCommand(path string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/C", "start", path)
	case "darwin":
		return exec.Command("open", path)
	default: // Linux
		return exec.Command("xdg-open", path)
	}
}

func SearchFolders(inputFolders []string, searchString string) (string, error) {
	start := time.Now()
	results := make(chan string, 100)
	done := make(chan bool)
	var wg sync.WaitGroup
	var matchedFolders []string

	go func() {
		fmt.Printf("\033[s") // save cursor position
		count := 0
		for result := range results {
			matchedFolders = append(matchedFolders, result)
			fmt.Printf("\033[u\033[J") // restore cursor position
			count++
		}
		done <- true
	}()

	for _, folder := range inputFolders {
		wg.Add(1)
		go func(folder string) {
			defer wg.Done()
			searchRecursive(ExpandPath(folder), searchString, &wg, results)
		}(folder)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	<-done // wait for all results to be printed

	elapsed := time.Since(start)
	fmt.Printf("%sOperation took %s%s\n", ColorGreen, elapsed, ColorReset)

	var index int
	if len(matchedFolders) == 0 {
		return "", fmt.Errorf("%sno matching folders found%s\n", ColorYellow, ColorReset)
	} else if len(matchedFolders) > 1 {
		fmt.Printf("%sMore than one project returned:\n", ColorYellow)
		for i, result := range matchedFolders {
			fmt.Printf("%s%d%s: %s\n", ColorBlue, i, ColorReset, result)
		}
		fmt.Printf("%sEnter index of selection: %s", ColorBoldGreen, ColorReset)

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return "", fmt.Errorf("error reading input")
		}

		response := scanner.Text()
		userIndex, err := strconv.Atoi(strings.TrimSpace(response))
		if err != nil {
			return "", fmt.Errorf("invalid selection")
		}
		if userIndex < 0 || userIndex >= len(matchedFolders) {
			return "", fmt.Errorf("invalid selection: out of range")
		}
		index = userIndex
	}
	return matchedFolders[index], nil
}

func searchRecursive(folderPath string, searchString string, wg *sync.WaitGroup, results chan string) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("Error reading folder %s: %v\n", folderPath, err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(folderPath, file.Name())

		// Check if folder matches search string
		if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(searchString)) {
			results <- fullPath
		}

		// Recursively search subdirectories
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			searchRecursive(path, searchString, wg, results)
		}(fullPath)
	}
}
