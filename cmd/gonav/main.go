package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var cfgFile string

// Config represents the structure of the YAML config file
type navConfig struct {
	DefaultFolders []string `yaml:"defaultFolders"`
	MaxDepth       int                  `yaml:"maxDepth"`
}

func expandPath(path string) string {
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


func main() {
    var rootCmd = &cobra.Command{
		Args:  cobra.NoArgs,
        Use:   "nav",
        Short: "Gonav is a CLI application",
        Long:  `Gonav is a CLI application written in Go.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
    }


	navCmd := (&cobra.Command{
		Use:   "go [folder]",
		Short: "Navigate to a project folder",
		Long: `Navigate to a project folder within the default folders specified in the configuration. 
		You can specify the depth of the search using the --depth flag. 
		Use the --code flag to open the selected folder with VS Code.`,
		Aliases: []string{"go"},
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			inputFolder := args[0]
			folders := viper.GetStringSlice("defaultFolders")
			if len(folders) == 0 {
				fmt.Println("No default folders found in the configuration.")
				return
			}

			var maxDepth int
			if flag := cmd.Flags().Lookup("depth"); flag != nil {
				maxDepth, _ = strconv.Atoi(flag.Value.String())
			} else {
				maxDepth = viper.GetInt("maxDepth")
			}
			fmt.Printf("Using max depth: %d\n", maxDepth)
			
			var wg sync.WaitGroup
			var results []string
			var mu sync.Mutex

			for _, folder := range folders {
				wg.Add(1)
				go func(folder string) {
					defer wg.Done()
					folder = expandPath(folder)
					files, err := os.ReadDir(folder)
					if err != nil {
						fmt.Printf("Error reading folder %s: %v\n", folder, err)
						return
					}

					for _, file := range files {
						if file.IsDir() && strings.Contains(strings.ToLower(file.Name()), strings.ToLower(inputFolder)) {
							mu.Lock()
							results = append(results, folder+"/"+file.Name())
							mu.Unlock()
						}
					}
				}(folder)
			}
			wg.Wait()

			if len(results) == 0 {
				fmt.Println("No matching folders found.")
			} else  if len(results) > 1{

				fmt.Println("More than one project returned:")
				for i, result := range results {
					fmt.Printf("%d: %s\n", i, result)
				}
				fmt.Println("Enter index of selection")
				var response string
				fmt.Scanln(&response)
				index, err := strconv.Atoi(response)
				if err != nil || index < 0 || index >= len(results) {
					fmt.Println("Invalid selection.")
					return
				}
				fmt.Printf("You selected: %s\n", results[index])
				code , _:= cmd.Flags().GetBool("code")
				if code {
					err = exec.Command("code", results[index]).Start()
					if err != nil {
						fmt.Println("Failed to open folder with code: %v\n", err)
					}
					return
				}

				
				err = exec.Command("xdg-open", results[index]).Start()
				if err != nil {
					fmt.Printf("Failed to open folder: %v\n", err)
				}
			}
		},
	})
	// Define flags
	navCmd.Flags().BoolP("code", "c", false, "Whether to open a folder with VS Code")
	
	rootCmd.AddCommand(navCmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:     "config",
		Short:   "Print the configuration",
		Aliases: []string{"pc"},
		Run: func(cmd *cobra.Command, args []string) {
			printConfig()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use: "add",
		Short: "Edit the config file",
		Aliases: []string{"a"},
		Run: func(cmd *cobra.Command, args []string) {
			configFile := viper.ConfigFileUsed()
	
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				fmt.Println("No config file found. Initialize with `gonav`")
				return
			}
			openInEditor(configFile)
		},
	})

    cobra.OnInitialize(initConfig)

    // rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gonav.yaml)")

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
func printConfig() {
	settings := viper.AllSettings()
	if len(settings) == 0 {
		fmt.Println("No configuration found.")
		os.Exit(1)
	}
	for key, value := range settings {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func createConfig(defaultConfigPath string) {
	fmt.Println("Config file does not exist. Do you want to create a default config file? (y/n)")
	var response string
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		fmt.Println("Creating default config file at " + defaultConfigPath)

		// In your main function:
		configYaml:= navConfig{
			DefaultFolders: []string{
				"~/Documents",
				"~/Projects",
			},
			MaxDepth: 3,
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
		openInEditor(defaultConfigPath)	
}}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        viper.AddConfigPath(home)
        viper.SetConfigName(".gonav")
		viper.SetConfigType("yaml")
    }

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err == nil {
		printConfigMessage(8, "/tmp/gonav_last_printed")
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

func printConfigMessage(hour int64, cacheFile string) {
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

func openInEditor(filePath string) {
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