// Copyright (c) Bruno Schaatsbergen
// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"fmt"
	"math"
	"net"
	"strconv"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/bschaatsbergen/cidr/pkg/helper"
	"github.com/spf13/cobra"
)

type DivideOptions struct {
	Network *net.IPNet
	Divisor int64
	Hosts   []int64
}

type ArgsKey struct{}

const (
	divideExample = "# Divides the given CIDR range into N distinct networks *Truncates output to 50\n" +
		"$ cidr divide 10.0.0.0/16 9\n" +
		"  [Networks]\n" +
		"10.0.0.0/20\n" +
		"10.0.16.0/20\n" +
		"10.0.32.0/20\n" +
		"10.0.48.0/20\n" +
		"10.0.64.0/20\n" +
		"10.0.80.0/20\n" +
		"10.0.96.0/20\n" +
		"10.0.112.0/20\n" +
		"10.0.128.0/20\n" +
		"-----\n" +
		"$ cidr divide 10.0.0/16 -u 23,10,125\n" +
		"  [Networks]            [Used]  [Total]\n" +
		"10.0.0.0/27               23      30\n" +
		"10.0.0.32/28              10      14\n" +
		"10.0.0.48/25             125     126\n"
)

var (
	divideCmd = &cobra.Command{
		Use:     "divide",
		Short:   "Divides CIDR range into a minimum of N distinct networks; -u flag for division by user count",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"d", "div"},
		Example: divideExample,
		PreRunE: divideArgumentValidator,
		RunE:    executeDivide,
	}
)

func init() {
	divideCmd.Flags().Int64SliceP("hosts", "H", []int64{}, "Divides by desired hosts instead of equal divisions. Say you want your network divided into 32 hosts, 10 hosts, and 17 hosts: 'cidr divide <network> 32,10,17' ")
	rootCmd.AddCommand(divideCmd)
}

// Validates the arguments for the network division operation.
// Catches network parsing errors and invalid divisors.
func divideArgumentValidator(cmd *cobra.Command, args []string) error {
	// Network parsing checks.
	network, err := core.ParseCIDR(args[0])
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	maskSize, _ := network.Mask.Size()
	if (helper.IsIPv4Network(network) && (maskSize == 32)) || maskSize >= 128 {
		return fmt.Errorf("Cannot divide a %s -- No Space\n", network.String())
	}

	// Divisor and -u Hosts validation.
	var divisor int64
	var desiredHosts []int64
	desiredHostsPerSubnet, err := cmd.Flags().GetInt64Slice("hosts")
	if err != nil {
		return err
	}
	if len(desiredHostsPerSubnet) == 0 {
		if len(args) < 2 {
			return fmt.Errorf("Enter at least one divisor\n")
		}
		divisor, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil || divisor <= 1 {
			return fmt.Errorf("%s\n", err)
		}
	} else {
		desiredHosts, err = core.ValidateHostSpace(network, desiredHostsPerSubnet)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}

	// Collects the valid arguments, so we don't need to bother with these checks later.
	validArgs := &DivideOptions{
		Network: network,
		Divisor: divisor,
		Hosts:   desiredHosts,
	}
	cmd.SetContext(context.WithValue(cmd.Context(), ArgsKey{}, validArgs))
	return nil
}

func executeDivide(cmd *cobra.Command, args []string) error {
	validArgsInterface := cmd.Context().Value(ArgsKey{})
	if validArgsInterface == nil {
		return fmt.Errorf("validArgs not in context\n")
	}
	validArgs, ok := validArgsInterface.(*DivideOptions)
	if !ok {
		return fmt.Errorf("Invalid type for validArgs: %T", validArgsInterface)
	}
	desiredHosts, err := cmd.Flags().GetInt64Slice("hosts")
	if err != nil {
		return err
	}
	if len(validArgs.Hosts) < 1 {
		networks, err := core.DivideCidr(validArgs.Network, validArgs.Divisor)
		if err != nil {
			return err
		}
		printNetworkPartitions(networks)
	} else {
		networks, err := core.DivideCidrHosts(validArgs.Network, validArgs.Hosts)
		if err != nil {
			return err
		}
		printNetworkHostBreakdown(networks, desiredHosts)

	}
	return nil
}


func printNetworkPartitions(networks []net.IPNet) {
	fmt.Println("  [Networks]")
	networkCount := len(networks)
	if networkCount <= 50 {
		for _, network := range networks {
			fmt.Println(network.String())
		}
	} else {
		for i := 0; i < networkCount; i++ {
			if i == 25 {
				i = networkCount - 25
				fmt.Println("......")
				continue
			}
			fmt.Println(networks[i].String())
		}
	}
}

func printNetworkHostBreakdown(networks []net.IPNet, desiredHosts []int64) {
	length := len(networks)
	fmt.Println("  [Networks]\t\t[Used]\t[Total]")
	for i := 0; i < length; i++ {
		hostMask, bits := core.GetNetmask(&networks[i]).Size()
		totalHosts := math.Pow(2, float64(bits-hostMask))
		fmt.Printf("%-18s%10d%8d\n", networks[i].String(), desiredHosts[i], int64(totalHosts-2.0))
	}
}

