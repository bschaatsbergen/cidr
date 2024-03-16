// Copyright (c) Bruno Schaatsbergen
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/spf13/cobra"
)

const (
	overlapsExample = "# Check whether 2 IPv4 CIDR ranges overlap\n" +
		"cidr overlaps 10.0.0.0/16 10.0.14.0/22\n" +
		"\n" +
		"# Check whether 2 IPv6 CIDR ranges overlap\n" +
		"cidr overlaps 2001:db8:1111:2222:1::/80 2001:db8:1111:2222:1:1::/96"
)

var overlapsCmd = &cobra.Command{
	Use:     "overlaps",
	Short:   "Checks if a CIDR range overlaps with another CIDR range",
	Example: overlapsExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("error: provide 2 CIDR ranges")
			fmt.Println("See 'cidr overlaps -h' for help and examples")
			os.Exit(1)
		}
		network1, err := core.ParseCIDR(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			fmt.Println("See 'cidr overlaps -h' for help and examples")
			os.Exit(1)
		}
		network2, err := core.ParseCIDR(args[1])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			fmt.Println("See 'cidr overlaps -h' for help and examples")
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
