package cmd

import (
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/bschaatsbergen/cidr/pkg/helper"
	"github.com/spf13/cobra"
)



const (
	divideExample = "# Divides the given CIDR range into N distinct networks\n" +
						 "cidr divide 10.0.0.0/16 9\n" +
						"[Networks]\n" +
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
						""
	divideHelpMessage = "See 'cidr divide -h' for more details"
);

var (
	divideCmd=&cobra.Command{
		Use: "divide",
		Short: "Divides CIDR range into N distinct ranges",
		Example: divideExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) !=2 {
				fmt.Printf("[ERROR]: Provide a CIDR range AND a divisor: cidr <CIDR_RANGE> <DIVISOR>\n")
				fmt.Println(divideHelpMessage)
				os.Exit(1)
			}
			network, err:= core.ParseCIDR(args[0])
			if err != nil {
				fmt.Printf("%s\n", err)
				fmt.Println(divideHelpMessage)
				os.Exit(1)
			}
			divisor, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || divisor <= 1 {
				fmt.Printf("%s\n", err)
				fmt.Println(divideHelpMessage)
				os.Exit(1)
			}
			networks, err := DivideCidr(network, divisor)
			if err != nil {
				fmt.Printf("%s\n", err)
				fmt.Println(divideHelpMessage)
				os.Exit(1)
			}
			printNetworkPartitions(networks)
		},

	}
)

func init() {
	rootCmd.AddCommand(divideCmd)
}

// Divides the given network into N smaller networks 
// Errors if division is not possible
func DivideCidr(network *net.IPNet, divisor uint64) ([]net.IPNet, error) {
	networks := make([]net.IPNet, 0)
	isIPv4 := helper.IsIPv4Network(network)
	maskSize, _ := network.Mask.Size()
	if divisor <= 1 {
		return nil, fmt.Errorf("[ERROR] Input a divisor higher than 1\n")
	}

	if isIPv4 && (maskSize == 32) {
		return nil, fmt.Errorf("[ERROR] Cannot divide a %s -- No Space\n",  network.String())
	}
	if maskSize >= 128 {
		return nil, fmt.Errorf("[ERROR] Cannot divide a %s -- No Space\n",  network.String())
	}

	addressCount := core.GetAddressCount(network)
	cidrWack, err := getPrefix(divisor, addressCount, isIPv4)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	nextAddress := new(net.IPNet) 
	nextAddress.IP = network.IP
	nextAddress.Mask = cidrWack
	for i := uint64(0); i < divisor ; i++  {
		networks = append(networks, *nextAddress)
		new, err := core.GetNextAddress(nextAddress, cidrWack) 
		if err != nil {
			return nil, err
		}
		nextAddress.IP = new.IP
	}
	return networks, nil
}

// Returns the net.IPMask necessary for the provided divisor. 
// Errors if address space if insufficient or divison is not possible.
func getPrefix(divisor uint64, addressCount *big.Int, IPv4 bool) (net.IPMask, error) {
	div := big.NewInt(int64(divisor))
	if addressCount.Cmp(div) == -1  || div.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("[ERROR] Cannot divide %d Addresses into %d divisions\n",addressCount, div)
	}
	
	// Gets address partitions 
	// If partition is 256
	//  
	// 00000010 
	// 2 << 8
	// 1 00000000
	// prefix = 8, then 32 - 8 = /24
	addressPartition := new(big.Int).Div(addressCount,div)
	two := big.NewInt(2)
	exponent := big.NewInt(0)
	for two.Cmp(addressPartition) <= 0 {
		two.Lsh(two,1)
		exponent.Add(exponent, big.NewInt(1))
	}
	subnetPrefix := int(exponent.Int64())
	if IPv4 {
		if subnetPrefix > 30 {
			return nil, fmt.Errorf("[ERROR] Address Space is insufficient for %d subnets\n", div)
		}
		return net.CIDRMask(32-subnetPrefix,32),nil
	}
	if subnetPrefix > 126 {
		return nil, fmt.Errorf("[ERROR] Address Space is insufficient for %d subnets\n", div)
	}
	return net.CIDRMask(128-subnetPrefix,128), nil

}


func printNetworkPartitions(networks []net.IPNet) {
	fmt.Println("[Networks]")
	for _, network := range networks {
		fmt.Println(network.String())
	}
}


