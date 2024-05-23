
package cmd

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/bschaatsbergen/cidr/pkg/helper"
	//"github.com/fatih/color"
	"github.com/spf13/cobra"
)



const (
	nextExample = "# Provides the next available address space with your desired CIDR\n" +
						 "cidr next 10.0.0.0/16 27\n" +
						 "Result: 10.1.0.0/27\n" 
	nextHelpMessage = "See 'cidr next -h' for more details"
);

var (
	nextCmd=&cobra.Command{
		Use: "next",
		Short: "Provides next available address space with your desired CIDR",
		Example: nextExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) !=2 {
				fmt.Printf("[ERROR]: Provide a network AND a cidr: cidr <CIDR_NETWORK> <CIDR>\n")
				fmt.Println(nextHelpMessage)
				os.Exit(1)
			}
			network, err:= core.ParseCIDR(args[0])
			if err != nil {
				fmt.Printf("%s\n", err)
				fmt.Println(nextHelpMessage)
				os.Exit(1)
			}
			cidrNum, err := strconv.ParseUint(args[1], 10, 64)
			var cidr net.IPMask
			if err != nil || cidrNum <= 1 {
				fmt.Printf("%s\n", err)
				fmt.Println(nextHelpMessage)
				os.Exit(1)
			}
			if helper.IsIPv4Network(network) {
				cidr = net.CIDRMask(int(cidrNum), 32)
			} else {
				cidr = net.CIDRMask(int(cidrNum), 128)
			}
			nextNetwork, err := core.GetNextAddress(network, cidr)
			if err != nil {
				fmt.Printf("%s\n", err)
				fmt.Println(nextHelpMessage)
				os.Exit(1)
			}
			fmt.Println(nextNetwork.String())
		},

	}
)

func init() {
	rootCmd.AddCommand(nextCmd)
}

