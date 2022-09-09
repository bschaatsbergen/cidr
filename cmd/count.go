package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/spf13/cobra"
)

var (
	countCmd = &cobra.Command{
		Use:   "count",
		Short: "Return the count of distinct host addresses in a given CIDR range",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("error: provide a CIDR range")
				fmt.Println("See 'cidr count -h' for help and examples")
				os.Exit(1)
			}
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				fmt.Printf("error: invalid CIDR range: %s\n", args[0])
				fmt.Println("See 'cidr count -h' for help and examples")
				os.Exit(1)
			}
			hostAddressCount := count(network)
			fmt.Println(hostAddressCount)

		},
	}
)

func init() {
	rootCmd.AddCommand(countCmd)
}

func count(network *net.IPNet) uint64 {
	count := core.AddressCount(network)
	return count
}
