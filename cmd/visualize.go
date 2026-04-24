package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/bschaatsbergen/cidr/internal/core"
	"github.com/spf13/cobra"
)

const (
	visualizeExample = "# Visualize an IPv4 CIDR range as a hierarchical tree\n" +
		"cidr visualize 10.0.0.0/29\n" +
		"\n" +
		"# Visualize with limited depth for large networks\n" +
		"cidr visualize 10.0.0.0/16 --depth 3\n" +
		"\n" +
		"# Show only summary information\n" +
		"cidr visualize 10.0.0.0/16 --compact"
)

var (
	maxDepth  int
	maxLines  int
	compact   bool
	noIPs     bool
	lineCount int

	visualizeCmd = &cobra.Command{
		Use:     "visualize",
		Short:   "Visualize a CIDR block as a hierarchical tree",
		Example: visualizeExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("error: provide a CIDR range")
				fmt.Println("See 'cidr visualize -h' for help and examples")
				os.Exit(1)
			}
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				fmt.Printf("error: %s\n", err)
				fmt.Println("See 'cidr visualize -h' for help and examples")
				os.Exit(1)
			}
			visualize(network)
		},
	}
)

func init() {
	rootCmd.AddCommand(visualizeCmd)
	visualizeCmd.Flags().IntVar(&maxDepth, "depth", 0, "Maximum recursion depth (0 = unlimited)")
	visualizeCmd.Flags().IntVar(&maxLines, "max-lines", 0, "Maximum output lines (0 = unlimited)")
	visualizeCmd.Flags().BoolVar(&compact, "compact", false, "Show compact summary instead of full tree")
	visualizeCmd.Flags().BoolVar(&noIPs, "no-ips", false, "Don't show individual IP addresses, only CIDR blocks")
}

func visualize(network *net.IPNet) {
	lineCount = 0
	fmt.Println(network.String())
	lineCount++

	if compact {
		visualizeCompact(network)
	} else {
		visualizeNetworkRecursive(network, "", false, 0)
		if maxLines > 0 && lineCount >= maxLines {
			fmt.Printf("\n... (output truncated at %d lines)\n", maxLines)
		}
	}
}

// visualizeCompact shows a compact summary of the CIDR block
func visualizeCompact(network *net.IPNet) {
	prefixLen, _ := network.Mask.Size()
	addressCount := core.GetAddressCount(network)

	fmt.Printf("├── Prefix Length: /%d\n", prefixLen)
	lineCount++

	if network.IP.To4() != nil {
		hostCount := core.GetHostAddressCount(network)
		fmt.Printf("├── Total Addresses: %s\n", addressCount.String())
		fmt.Printf("├── Usable Hosts: %s\n", hostCount.String())
	} else {
		fmt.Printf("├── Total Addresses: %s\n", addressCount.String())
	}

	// Calculate last address
	lastIP := make(net.IP, len(network.IP))
	for i := range lastIP {
		lastIP[i] = network.IP[i] | ^network.Mask[i]
	}

	fmt.Printf("├── First Address: %s\n", network.IP.String())
	fmt.Printf("└── Last Address: %s\n", lastIP.String())
	lineCount += 4
}

// visualizeNetworkRecursive recursively visualizes the network structure
func visualizeNetworkRecursive(network *net.IPNet, prefix string, isLast bool, depth int) {
	// Check line limit
	if maxLines > 0 && lineCount >= maxLines {
		return
	}

	// Check depth limit
	if maxDepth > 0 && depth >= maxDepth {
		return
	}

	prefixLen, bits := network.Mask.Size()

	// For IPv4, if we're at a /30 or smaller, show the individual IPs (unless --no-ips)
	if network.IP.To4() != nil && prefixLen >= 30 && !noIPs {
		showSubnetsWithPrefix(network, prefix, isLast)
		return
	}

	// For IPv6, if we're at a /126 or smaller, show the individual IPs (unless --no-ips)
	if network.IP.To16() != nil && network.IP.To4() == nil && prefixLen >= 126 && !noIPs {
		showSubnetsWithPrefix(network, prefix, isLast)
		return
	}

	// Calculate next level subnets (increase prefix by 1)
	nextPrefixLen := prefixLen + 1
	subnets := getSubnets(network, nextPrefixLen)

	for i, subnet := range subnets {
		// Check line limit before printing
		if maxLines > 0 && lineCount >= maxLines {
			return
		}

		isLastSubnet := (i == len(subnets)-1)

		// Create the tree characters
		var connector string
		var newPrefix string
		if isLastSubnet {
			connector = "└── "
			newPrefix = prefix + "    "
		} else {
			connector = "├── "
			newPrefix = prefix + "│   "
		}

		fmt.Println(prefix + connector + subnet.String())
		lineCount++

		// Recurse if we haven't reached the leaf level
		if nextPrefixLen < bits-1 {
			visualizeNetworkRecursive(subnet, newPrefix, isLastSubnet, depth+1)
		} else if nextPrefixLen <= bits && !noIPs {
			// Show final level subnets
			showSubnetsWithPrefix(subnet, newPrefix, isLastSubnet)
		}
	}
}

// showSubnetsWithPrefix shows individual IPs with proper prefix
func showSubnetsWithPrefix(network *net.IPNet, prefix string, parentIsLast bool) {
	prefixLen, bits := network.Mask.Size()

	// Limit the number of addresses to display
	// Calculate the number of bits available for host addresses
	hostBits := bits - prefixLen
	if hostBits > 8 {
		// Too many addresses to display individually (more than 256)
		return
	}

	// Get total number of addresses
	addressCount := uint64(1 << hostBits)

	ips := make([]net.IP, 0)
	ip := net.ParseIP(network.IP.String())

	for i := uint64(0); i < addressCount; i++ {
		// Check line limit
		if maxLines > 0 && lineCount >= maxLines {
			break
		}

		ips = append(ips, net.ParseIP(ip.String()))
		incrementIP(ip)
	}

	for i, ip := range ips {
		// Check line limit
		if maxLines > 0 && lineCount >= maxLines {
			return
		}

		isLast := (i == len(ips)-1)

		var connector string
		if isLast {
			connector = "└── "
		} else {
			connector = "├── "
		}

		fmt.Println(prefix + connector + ip.String())
		lineCount++
	}
}

// incrementIP increments the IP address by 1
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			return
		}
	}
}

// getSubnets returns the two subnets by splitting the given network
func getSubnets(network *net.IPNet, newPrefixLen int) []*net.IPNet {
	prefixLen, _ := network.Mask.Size()
	if newPrefixLen <= prefixLen {
		return []*net.IPNet{network}
	}

	// Create two subnets by splitting at the new prefix length
	subnet1 := &net.IPNet{
		IP:   network.IP,
		Mask: net.CIDRMask(newPrefixLen, 32),
	}
	if network.IP.To4() == nil {
		subnet1.Mask = net.CIDRMask(newPrefixLen, 128)
	}

	// Calculate second subnet
	subnet2IP := incrementIPByPowerOfTwo(network.IP, newPrefixLen)
	subnet2 := &net.IPNet{
		IP:   subnet2IP,
		Mask: subnet1.Mask,
	}

	return []*net.IPNet{subnet1, subnet2}
}

// incrementIPByPowerOfTwo increments the IP by 2^(bits - prefixLen)
func incrementIPByPowerOfTwo(ip net.IP, prefixLen int) net.IP {
	result := make(net.IP, len(ip))
	copy(result, ip)

	bits := 32
	if ip.To4() == nil {
		bits = 128
	}

	increment := uint32(1 << (bits - prefixLen))

	// Convert to integer, add increment, convert back
	if ip.To4() != nil {
		// IPv4
		ipUint := uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
		ipUint += increment
		result[0] = byte((ipUint >> 24) & 0xFF)
		result[1] = byte((ipUint >> 16) & 0xFF)
		result[2] = byte((ipUint >> 8) & 0xFF)
		result[3] = byte(ipUint & 0xFF)
	}

	return result
}
