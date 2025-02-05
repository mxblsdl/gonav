package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var cfgFile string

// CommentedStringSlice represents a slice of strings with header comment
type CommentedStringSlice struct {
    Values  []string `yaml:"values"`
    Comment string   `yaml:"comment,omitempty"`
}

// Config represents the structure of the YAML config file
type Yaml struct {
	DefaultFolders CommentedStringSlice `yaml:"defaultFolders"`
    MaxDepth       CommentedStringSlice `yaml:"maxDepth"`
}

// MarshalYAML implements custom marshaling for CommentedStringSlice
func (c CommentedStringSlice) MarshalYAML() (interface{}, error) {
    var result string
    if c.Comment != "" {
        result = "# " + c.Comment + "\n"
    }
    
    for _, v := range c.Values {
        result += "- " + v + "\n"
    }
    
    return result, nil
}

func main() {
    var rootCmd = &cobra.Command{
		Args:  cobra.NoArgs,
        Use:   "gonav",
        Short: "Gonav is a CLI application",
        Long:  `Gonav is a CLI application written in Go.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
    }

	rootCmd.AddCommand(&cobra.Command{
		Use: "nav",
		Short: "Navigate to a project folder",
		Aliases: []string{"nav"},
		Run: func(cmd *cobra.Command, args[]string) {

		},
	})

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

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gonav.yaml)")

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
	fmt.Println(settings)
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
		configYaml := Yaml{
			DefaultFolders: CommentedStringSlice{
				Values: []string{
					"~/Documents",
					"~/Projects",
				},
				Comment: "List of default folders to search",
			},
			MaxDepth: CommentedStringSlice{
				Values:   []string{"3"},
				Comment: "Maximum depth to search in directories",
			},
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
        fmt.Println("Config file found:", viper.ConfigFileUsed())
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