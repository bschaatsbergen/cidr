package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/bschaatsbergen/cidr/pkg/core"
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
			broadcastAddress, netmask, prefixLength, baseAdress, count, firstUsableIPAddress, lastUsableIPAddress := getNetworkDetails(network)
			explain(broadcastAddress, netmask, prefixLength, baseAdress, count, firstUsableIPAddress, lastUsableIPAddress)
		},
	}
)

func init() {
	rootCmd.AddCommand(explainCmd)
}

func getNetworkDetails(network *net.IPNet) (string, string, string, string, string, string, string) {
	// Broadcast address
	var broadcastAddress string

	ipBroadcast, err := core.GetBroadcastAddress(network)
	if err != nil {
		broadcastAddress = err.Error()
	} else {
		broadcastAddress = ipBroadcast.String()
	}

	// Netmask
	netmask := core.GetNetMask(network)
	prefixLength := core.GetPrefixLength(netmask)

	// Base address
	baseAddress := core.GetBaseAddress(network)

	// Address count
	var count string
	addressCount := core.GetAddressCount(network)
	// Produce a human readable number
	count = message.NewPrinter(language.English).Sprintf("%d", addressCount)

	// First usable IP address
	var firstUsableIPAddress string
	firstUsableIP, err := core.GetFirstUsableIPAddress(network)
	if err != nil {
		firstUsableIPAddress = err.Error()
	} else {
		firstUsableIPAddress = firstUsableIP.String()
	}

	// Last usable IP address
	var lastUsableIPAddress string
	lastUsableIP, err := core.GetLastUsableIPAddress(network)
	if err != nil {
		lastUsableIPAddress = err.Error()
	} else {
		lastUsableIPAddress = lastUsableIP.String()
	}

	return broadcastAddress, netmask.String(), fmt.Sprint(prefixLength), baseAddress.String(), count, firstUsableIPAddress, lastUsableIPAddress
}

//nolint:goconst
func explain(broadcastAddress, netmask, prefixLength, baseAddress, count, firstUsableIPAddress, lastUsableIPAddress string) {
	fmt.Printf(color.BlueString("Base Address:\t\t ")+"%s\n", baseAddress)
	fmt.Printf(color.BlueString("Usable IP Address range: ")+"%s to %s\n", firstUsableIPAddress, lastUsableIPAddress)
	fmt.Printf(color.BlueString("Broadcast Address:\t ")+"%s\n", broadcastAddress)
	fmt.Printf(color.BlueString("Address Count:\t\t ")+"%s\n", count)
	fmt.Printf(color.BlueString("Netmask:\t\t ")+"%s (/%s bits)\n", netmask, prefixLength)
}
