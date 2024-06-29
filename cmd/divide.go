package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/bschaatsbergen/cidr/pkg/helper"
	"github.com/spf13/cobra"
)

const (
	divideExample = "# Divides the given CIDR range into N distinct networks *Truncates output to 50\n" +
		"$ cidr divide 10.0.0.0/16 9\n" +
		"10.0.0.0/20\n" +
		"10.0.16.0/20\n" +
		"10.0.32.0/20\n" +
		"10.0.48.0/20\n" +
		"10.0.64.0/20\n" +
		"10.0.80.0/20\n" +
		"10.0.96.0/20\n" +
		"10.0.112.0/20\n" +
		"10.0.128.0/20\n"
)

var divideCmd = &cobra.Command{
	Use:     "divide",
	Short:   "Divides the given CIDR range into N distinct networks",
	Args:    cobra.MinimumNArgs(2),
	Example: divideExample,
	PreRunE: validateDivideArguments,
	RunE:    executeDivide,
}

func init() {
	rootCmd.AddCommand(divideCmd)
}

func validateDivideArguments(cmd *cobra.Command, args []string) error {
	// Ensure CIDR is valid
	_, _, err := net.ParseCIDR(args[0])
	if err != nil {
		return fmt.Errorf("invalid network: %s", args[0])
	}

	// Ensure divisor is a valid integer
	_, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid divisor: %s", args[1])
	}

	return nil
}

func executeDivide(cmd *cobra.Command, args []string) error {
	network, err := core.ParseCIDR(args[0])
	if err != nil {
		return fmt.Errorf("invalid network: %s", args[0])
	}

	maskSize, _ := network.Mask.Size()
	if (helper.IsIPv4Network(network) && maskSize == 32) || maskSize >= 128 {
		return fmt.Errorf("invalid network mask size: %s", args[0])
	}

	divisor, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid divisor: %s", args[1])
	}

	networks, err := core.DivideCIDR(network, divisor)
	if err != nil {
		return err
	}

	printNetworkPartitions(networks)
	return nil
}

func printNetworkPartitions(networks []net.IPNet) {
	const truncateLimit = 50
	networkCount := len(networks)

	if networkCount <= truncateLimit {
		for _, network := range networks {
			fmt.Println(network.String())
		}
	} else {
		for i := 0; i < truncateLimit/2; i++ {
			fmt.Println(networks[i].String())
		}
		fmt.Println("......")
		for i := networkCount - truncateLimit/2; i < networkCount; i++ {
			fmt.Println(networks[i].String())
		}
	}
}
