package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mxblsdl/gonav/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var openExplorer bool

var navCmd = (&cobra.Command{
	Use:   "go [folder]",
	Short: "Navigate to a project folder",
	Long: `Navigate to a project folder within the default folders specified in the configuration.
	By default, launch an interactive shell in the selected folder.
	Use --open to open the selected folder in your file explorer instead.`,
	Aliases: []string{"go"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchTerm := args[0]
		folders := viper.GetStringSlice("Folders")
		if len(folders) == 0 {
			return fmt.Errorf("no default folders found in the configuration")
		}

		matchFolder, err := helpers.SearchFolders(folders, searchTerm)
		if err != nil {
			return err
		}

		if openExplorer {
			command := helpers.OpenShellCommand(matchFolder)
			if err := command.Start(); err != nil {
				return fmt.Errorf("error opening folder: %w", err)
			}
			return nil
		}

		return helpers.RunShellInDirectory(matchFolder)
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
		// Try to find a .code-workspace file in the folder
		workspaceFile := filepath.Join(matchFolder, "*.code-workspace")
		files, _ := filepath.Glob(workspaceFile)
		var command *exec.Cmd

		if len(files) > 0 {
			// Open the workspace file
			command = exec.Command("code", files[0])
			fmt.Printf("Running command: %s%s%s\n", helpers.BgGreen, command, helpers.ColorReset)
		} else {
			// Fallback to opening as a folder
			command = exec.Command("code", matchFolder)
			fmt.Printf("Running command: %s%s%s", helpers.BgGreen, command, helpers.ColorReset)
		}

		output, err := command.CombinedOutput()
		if err != nil {
			fmt.Printf("Error occurred while running command %q: %v\n", command.String(), err)
			if len(output) > 0 {
				fmt.Printf("Command output: %s", output)
			}
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
		output, err := command.CombinedOutput()
		if err != nil {
			fmt.Printf("Error occurred while running command %q: %v\n", command.String(), err)
			if len(output) > 0 {
				fmt.Printf("Command output: %s", output)
			}
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
	navCmd.Flags().BoolVarP(&openExplorer, "open", "o", false, "Open the selected folder in the file explorer")
	rootCmd.AddCommand(navCmd)
	rootCmd.AddCommand(printCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(codeCmd)
	rootCmd.AddCommand(posiCmd)
}
