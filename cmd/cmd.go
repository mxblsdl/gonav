package cmd

import (
	"fmt"
	"os"
	"os/exec"

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
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchTerm := args[0]
		folders := viper.GetStringSlice("Folders")
		if len(folders) == 0 {
			fmt.Printf("%sNo default folders found in the configuration.", helpers.ColorBoldRed)
			return
		}

		matchFolder, err := helpers.SearchFolders(folders, searchTerm)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(0)
		}

		fmt.Printf("You selected: %s\n%s", matchFolder, helpers.ColorReset)
		command := helpers.OpenShellCommand(matchFolder)
		err = command.Start()
		if err != nil {
			fmt.Println("Error opening folder:", err)
			os.Exit(1)
		}

	},
})

var codeCmd = (&cobra.Command{
	Use:     "code [folder]",
	Short:   "Open a folder with VS Code",
	Aliases: []string{"c"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchTerm := args[0]
		folders := viper.GetStringSlice("Folders")
		if len(folders) == 0 {
			fmt.Printf("%sNo default folders found in the configuration.", helpers.ColorBoldRed)
			return
		}

		matchFolder, err := helpers.SearchFolders(folders, searchTerm)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(0)
		}

		fmt.Printf("You selected: %s\n%s", matchFolder, helpers.ColorReset)
		command := exec.Command("code", matchFolder)
		err = command.Start()
		if err != nil {
			fmt.Println("Error opening folder:", err)
			os.Exit(1)
		}
	},
})

var posiCmd = (&cobra.Command{
	Use:     "positron [folder]",
	Short:   "Open a folder with Positron",
	Aliases: []string{"p"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchTerm := args[0]
		folders := viper.GetStringSlice("Folders")
		if len(folders) == 0 {
			fmt.Printf("%sNo default folders found in the configuration.", helpers.ColorBoldRed)
			return
		}

		matchFolder, err := helpers.SearchFolders(folders, searchTerm)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(0)
		}

		fmt.Printf("You selected: %s\n%s", matchFolder, helpers.ColorReset)
		command := exec.Command("positron", matchFolder)
		err = command.Start()
		if err != nil {
			fmt.Println("Error opening folder:", err)
			os.Exit(1)
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
			fmt.Printf("%sNo configuration found.\n", helpers.ColorRed)
			os.Exit(1)
		}
		for key, value := range settings {
			fmt.Printf("%s%s%s: %v\n", helpers.ColorBlue, key, helpers.ColorReset, value)
		}
	},
})

var addCmd = (&cobra.Command{
	Use:     "add",
	Short:   "Edit the config file",
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		configFile := viper.ConfigFileUsed()

		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			fmt.Printf("%sNo config file found. Initialize with `gonav`", helpers.ColorRed)
			return
		}
		err := helpers.OpenInEditor(configFile)
		if err != nil {
			fmt.Printf("%sError opening config file: %v%s\n", helpers.ColorRed, err, helpers.ColorReset)
			return
		}
	},
})

func init() {
	rootCmd.AddCommand(navCmd)
	rootCmd.AddCommand(printCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(codeCmd)
	rootCmd.AddCommand(posiCmd)
}
