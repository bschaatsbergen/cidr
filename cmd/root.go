package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version string

	rootCmd = &cobra.Command{
		Use:     "cidr",
		Short:   "cidr - CLI to perform various actions on CIDR ranges",
		Version: version, // The version is set during the build by making using of `go build -ldflags`
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
