package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/mxblsdl/gonav/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var navCmd = (&cobra.Command{
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
				folder = helpers.ExpandPath(folder)
				// TODO make recursive here with Walk
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
					fmt.Printf("Failed to open folder with code: %v\n", err)
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

var printCmd = (&cobra.Command{
	Use:     "config",
	Short:   "Print the configuration",
	Aliases: []string{"pc"},
	Run: func(cmd *cobra.Command, args []string) {
			settings := viper.AllSettings()
			if len(settings) == 0 {
				fmt.Println("No configuration found.")
				os.Exit(1)
			}
			for key, value := range settings {
				fmt.Printf("%s: %v\n", key, value)
			}
	},
})


var addCmd = (&cobra.Command{
	Use: "add",
	Short: "Edit the config file",
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		configFile := viper.ConfigFileUsed()

		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			fmt.Println("No config file found. Initialize with `gonav`")
			return
		}
		helpers.OpenInEditor(configFile)
	},
})

func init() {
	rootCmd.AddCommand(navCmd)
	rootCmd.AddCommand(printCmd)
	rootCmd.AddCommand(addCmd)
		
	// Define flags
	navCmd.Flags().BoolP("code", "c", false, "Whether to open a folder with VS Code")
}
