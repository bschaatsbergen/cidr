package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/spf13/cobra"
)

var overlapsCmd = &cobra.Command{
	Use:   "overlaps",
	Short: "Checks if a CIDR range overlaps with another CIDR range",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Provide 2 cidr ranges")
			os.Exit(1)
		}
		network1, err := core.ParseCIDR(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		network2, err := core.ParseCIDR(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		networksOverlap := overlaps(network1, network2)
		fmt.Println(networksOverlap)
	},
}

func init() {
	rootCmd.AddCommand(overlapsCmd)
}

func overlaps(network1, network2 *net.IPNet) bool {
	overlaps := core.Overlaps(network1, network2)
	return overlaps
}
