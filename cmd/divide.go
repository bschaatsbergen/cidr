// Copyright (c) Bruno Schaatsbergen
// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"net"
	"strconv"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/bschaatsbergen/cidr/pkg/helper"
	"github.com/spf13/cobra"
)

type DivideOptions struct {
	Network *net.IPNet
	Divisor int64
	Users   []int64
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
	divideCmd.Flags().Int64SliceP("users", "u", []int64{}, "Divides by desired users instead of equal divisions. Say you want your network divided into 32 users, 10 users, and 17 users: 'cidr divide <network> 32,10,17' ")
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
	if helper.IsIPv4Network(network) && (maskSize == 32) {
		return fmt.Errorf("Cannot divide a %s -- No Space\n", network.String())
	}
	if maskSize >= 128 {
		return fmt.Errorf("Cannot divide a %s -- No Space\n", network.String())
	}

	// Divisor and -u Users validation.
	var divisor int64
	var desiredUsers []int64
	desiredUsersPerSubnet, err := cmd.Flags().GetInt64Slice("users")
	if err != nil {
		return err
	}
	if len(desiredUsersPerSubnet) == 0 {
		if len(args) < 2 {
			return fmt.Errorf("Enter at least one divisor\n")
		}
		div, err := strconv.ParseInt(args[1], 10, 64)
		divisor = div
		if err != nil || div <= 1 {
			return fmt.Errorf("%s\n", err)
		}
	} else {
		desiredUsers, err = ValidateUserSpace(network, desiredUsersPerSubnet)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}

	// Collects the valid arguments, so we don't need to bother with these checks later.
	validArgs := &DivideOptions{
		Network: network,
		Divisor: divisor,
		Users:   desiredUsers,
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
	desiredUsers, err := cmd.Flags().GetInt64Slice("users")
	if err != nil {
		return err
	}
	if len(validArgs.Users) < 1 {
		networks, err := DivideCidr(validArgs.Network, validArgs.Divisor)
		if err != nil {
			return err
		}
		printNetworkPartitions(networks)
	} else {
		networks, err := DivideCidrHosts(validArgs.Network, validArgs.Users)
		if err != nil {
			return err
		}
		printNetworkUserBreakdown(networks, desiredUsers)

	}
	return nil
}

// Divides the given network into N smaller networks.
// Errors if division is not possible.
func DivideCidr(network *net.IPNet, divisor int64) ([]net.IPNet, error) {
	isIPv4 := helper.IsIPv4Network(network)

	addressCount := core.GetAddressCount(network)
	cidrWack, err := getPrefix(divisor, addressCount, isIPv4)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	networks := make([]net.IPNet, divisor)
	nextAddress := new(net.IPNet)
	nextAddress.IP = network.IP
	nextAddress.Mask = cidrWack
	subnetSize := core.GetAddressCount(nextAddress)
	for i := int64(0); i < divisor; i++ {
		networks[i] = *nextAddress
		ipAsInt := new(big.Int).SetBytes(nextAddress.IP)
		nextAddress.IP = new(big.Int).Add(ipAsInt, subnetSize).Bytes()
	}
	return networks, nil
}

// Given an IP Network and a list of desired users, splits the network into a configuration that accommodates the user space.
// Eg. Splitting a 192.168.0.0/24 where I need a network with 20 users, 10 users, and 30users.
// Gets the /27 /28 and /27 combination.
func DivideCidrHosts(network *net.IPNet, desiredUsers []int64) ([]net.IPNet, error) {

	var bits int
	if helper.IsIPv4Network(network) {
		bits = net.IPv4len * 8
	} else {
		bits = net.IPv6len * 8
	}
	networks := make([]net.IPNet, len(desiredUsers))
	nextIP := network.IP
	for i, users := range desiredUsers {
		subnetPrefix := net.CIDRMask(bits-int(users), bits)
		nextNetwork := net.IPNet{
			IP:   nextIP,
			Mask: subnetPrefix,
		}
		networks[i] = nextNetwork
		nextNetwork, err := core.GetNextAddress(&nextNetwork, subnetPrefix)
		if err != nil {
			return nil, err
		}
		nextIP = nextNetwork.IP
	}
	return networks, nil
}

// Returns the net.IPMask necessary for the provided divisor.
// Errors if address space if insufficient or division is not possible.
func getPrefix(divisor int64, addressCount *big.Int, IPv4 bool) (net.IPMask, error) {
	div := big.NewInt(divisor)
	if addressCount.Cmp(div) == -1 || div.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("Cannot divide %d Addresses into %d divisions\n", addressCount, div)
	}

	addressPartition := new(big.Int).Div(addressCount, div)
	two := big.NewInt(2)
	exponent := big.NewInt(0)
	for two.Cmp(addressPartition) <= 0 {
		two.Lsh(two, 1)
		exponent.Add(exponent, big.NewInt(1))
	}
	subnetPrefix := int(exponent.Int64())
	bits := net.IPv6len * 8
	if IPv4 {
		bits = net.IPv4len * 8
		if subnetPrefix > 30 {
			return nil, fmt.Errorf("Address Space is insufficient for %d subnets\n", div)
		}
	}
	if subnetPrefix > 126 {
		return nil, fmt.Errorf("Address Space is insufficient for %d subnets\n", div)
	}
	return net.CIDRMask(bits-subnetPrefix, bits), nil

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

func printNetworkUserBreakdown(networks []net.IPNet, desiredUsers []int64) {
	length := len(networks)
	fmt.Println("  [Networks]\t\t[Used]\t[Total]")
	for i := 0; i < length; i++ {
		hostMask, bits := core.GetNetmask(&networks[i]).Size()
		totalHosts := math.Pow(2, float64(bits-hostMask))
		fmt.Printf("%-18s%10d%8d\n", networks[i].String(), desiredUsers[i], int64(totalHosts-2.0))
	}
}

func ValidateUserSpace(network *net.IPNet, desiredUsersPerSubnet []int64) ([]int64, error) {
	requiredAddressSpace := int64(0)
	totalAddressSpace := core.GetAddressCount(network)
	desiredUsers := make([]int64, len(desiredUsersPerSubnet))
	for i, users := range desiredUsersPerSubnet {
		subnetExponent := math.Ceil(math.Log2(float64(users) + 2.0))
		addressSpace := int64(math.Pow(2, subnetExponent))
		requiredAddressSpace += addressSpace
		desiredUsers[i] = int64(subnetExponent)
	}
	requiredAddressSpaceBig := big.NewInt(requiredAddressSpace)
	if totalAddressSpace.Cmp(requiredAddressSpaceBig) < 0 {
		return nil, fmt.Errorf("Total address space is: %s but desired Users requires %d addresses\n", totalAddressSpace.String(), requiredAddressSpace)
	}
	return desiredUsers, nil
}
