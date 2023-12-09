package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/bschaatsbergen/cidr/pkg/helper"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
			details := getNetworkDetails(network)
			explain(details)
		},
	}
)

func init() {
	rootCmd.AddCommand(explainCmd)
}

type networkDetailsToDisplay struct {
	IsIPV4Network              bool
	IsIPV6Network              bool
	BroadcastAddress           string
	BroadcastAddressHasError   bool
	Netmask                    net.IP
	PrefixLength               int
	BaseAddress                net.IP
	Count                      string
	UsableAddressRangeHasError bool
	FirstUsableIPAddress       string
	LastUsableIPAddress        string
}

func getNetworkDetails(network *net.IPNet) *networkDetailsToDisplay {
	details := &networkDetailsToDisplay{}

	// Determine whether the network is an IPv4 or IPv6 network.
	if helper.IsIPv4Network(network) {
		details.IsIPV4Network = true
	} else if helper.IsIPv6Network(network) {
		details.IsIPV6Network = true
	}

	// Obtain the broadcast address, handling errors if they occur.
	ipBroadcast, err := core.GetBroadcastAddress(network)
	if err != nil {
		// Set error flags and store the error message so that it can be displayed later.
		details.BroadcastAddressHasError = true
		details.BroadcastAddress = err.Error()
	} else {
		details.BroadcastAddress = ipBroadcast.String()
	}

	// Obtain the netmask and prefix length.
	netmask := core.GetNetmask(network)
	// A human-readable representation of the netmask is displayed in the output.
	details.Netmask = core.NetMaskToIPAddress(netmask)
	details.PrefixLength = core.GetPrefixLength(details.Netmask)

	// Obtain the base address of the network.
	details.BaseAddress = core.GetBaseAddress(network)

	// Obtain the total count of addresses in the network.
	count := core.GetAddressCount(network)
	// Format the count as a human-readable string and store it in the details struct.
	details.Count = message.NewPrinter(language.English).Sprintf("%d", count)

	// Obtain the first and last usable IP addresses, handling errors if they occur.
	firstUsableIP, err := core.GetFirstUsableIPAddress(network)
	if err != nil {
		// Set error flags if an error occurs during the retrieval of the first usable IP address.
		details.UsableAddressRangeHasError = true
	} else {
		details.FirstUsableIPAddress = firstUsableIP.String()
	}

	lastUsableIP, err := core.GetLastUsableIPAddress(network)
	if err != nil {
		// Set error flags if an error occurs during the retrieval of the last usable IP address.
		details.UsableAddressRangeHasError = true
	} else {
		details.LastUsableIPAddress = lastUsableIP.String()
	}

	// Return the populated 'networkDetailsToDisplay' struct.
	return details
}

//nolint:goconst
func explain(details *networkDetailsToDisplay) {
	var lengthIndicator string

	fmt.Printf(color.BlueString("Base Address:\t\t ")+"%s\n", details.BaseAddress)
	if !details.UsableAddressRangeHasError {
		fmt.Printf(color.BlueString("Usable Address Range:\t ")+"%s to %s\n", details.FirstUsableIPAddress, details.LastUsableIPAddress)
	} else {
		fmt.Printf(color.RedString("Usable Address Range:\t ")+"%s\n", "unable to calculate usable address range")
	}
	if !details.BroadcastAddressHasError && details.IsIPV4Network {
		fmt.Printf(color.BlueString("Broadcast Address:\t ")+"%s\n", details.BroadcastAddress)
	} else if details.BroadcastAddressHasError && details.IsIPV4Network {
		fmt.Printf(color.RedString("Broadcast Address:\t ")+"%s\n", details.BroadcastAddress)
	}
	fmt.Printf(color.BlueString("Address Count:\t\t ")+"%s\n", details.Count)

	if details.PrefixLength > 1 {
		lengthIndicator = "bits"
	} else {
		lengthIndicator = "bit"
	}

	fmt.Printf(color.BlueString("Netmask:\t\t ")+"%s (/%d %s)\n", details.Netmask, details.PrefixLength, lengthIndicator)
}
