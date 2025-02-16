package cmd

import (
	"fmt"
	"os"

	"github.com/mxblsdl/gonav/helpers"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "nav",
	Short: "Gonav is a CLI application",
	Long:  `Gonav is a CLI application written in Go.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}


func Execute() {
	cobra.OnInitialize(helpers.InitConfig)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

