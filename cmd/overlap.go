package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var overlapCmd = &cobra.Command{
	Use:   "overlap",
	Short: "Checks whether 2 cidr ranges overlap",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("overlap called")
	},
}

func init() {
	rootCmd.AddCommand(overlapCmd)
}
