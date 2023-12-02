package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/spf13/cobra"
)

var (
	explainCmd = &cobra.Command{
		Use:   "explain",
		Short: "Provides information about a CIDR range",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("error: provide a CIDR range and an IP address")
				fmt.Println("See 'cidr contains -h' for help and examples")
				os.Exit(1)
			}
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				fmt.Printf("error: %s\n", err)
				fmt.Println("See 'cidr contains -h' for help and examples")
				os.Exit(1)
			}
			foo, bar := explain(network)
			fmt.Println(foo)
			fmt.Println(bar)
		},
	}
)

func init() {
	rootCmd.AddCommand(explainCmd)
}

func explain(network *net.IPNet) (string, string) {
	return "", ""
}
