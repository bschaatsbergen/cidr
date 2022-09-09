package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/spf13/cobra"
)

var (
	containsCmd = &cobra.Command{
		Use:   "contains",
		Short: "Checks whether an IP address belongs to a CIDR range",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Println("Provide an IP address and a CIDR range")
				os.Exit(1)
			}
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			ip := core.ParseIP(args[1])
			if ip == nil {
				fmt.Println("Provided ip is not valid")
				os.Exit(1)
			}
			networkContainsIP := contains(network, ip)
			fmt.Println(networkContainsIP)
		},
	}
)

func init() {
	rootCmd.AddCommand(containsCmd)
}

func contains(network *net.IPNet, ip net.IP) bool {
	contains := core.ContainsAddress(network, ip)
	return contains
}
