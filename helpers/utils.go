package helpers

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func OpenInEditor(filePath string) {
    editor := os.Getenv("EDITOR")
    if editor == "" {
        editor = "nano" // default to nano if EDITOR is not set
    }
    cmd := exec.Command(editor, filePath)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error opening config file in editor:", err)
        os.Exit(1)
    }
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
		os.WriteFile(cacheFile, []byte(time.Now().Format(time.RFC3339)), 0644)
	}
}

func createConfig(defaultConfigPath string) {
	fmt.Printf("Config file does not exist. Do you want to create a default config file? %s(y/n)%s\n", ColorGreen, ColorReset)
	var response string
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		fmt.Println("Creating default config file at " + defaultConfigPath)

		configYaml:= navConfig{
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
		OpenInEditor(defaultConfigPath)	
}}

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
		PrintConfigMessage(8, "/tmp/gonav_last_printed")
    } else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defaultConfigPath := home + "/.gonav.yaml"
		createConfig(defaultConfigPath)
	}
}

func ScanWithWalkDir(root string) ([]string,error) {
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