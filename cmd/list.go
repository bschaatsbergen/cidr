package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/spf13/cobra"
)

const (
	listExample = "# Return a list of all ip addresses within a given IPv4 CIDR range\n" +
		"cidr list 10.0.0.0/28\n" +
		"\n" +
		"# Return the count of all distinct host addresses within a given IPv6 CIDR range\n" +
		"cidr list 2001:db8:1234:1a00::/126\n"
)

var (
	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "Return the list of all matching addresses in a given CIDR range",
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("error: provide a CIDR range")
				fmt.Println("See 'cidr list -h' for help and examples")
				os.Exit(1)
			}
			err := core.ListCIDR(args[0])
			if err != nil {
				fmt.Printf("error: invalid CIDR range: %s\n", args[0])
				fmt.Println("See 'cidr list -h' for help and examples")
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func list(network *net.IPNet) uint64 {
	count := core.AddressCount(network)
	return count
}
