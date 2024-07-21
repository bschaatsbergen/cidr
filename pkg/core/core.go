package core

import (
	"errors"
	"fmt"
	"math/big"
	"net"

	"github.com/bschaatsbergen/cidr/pkg/helper"
)

// ParseCIDR parses the given CIDR notation string and returns the corresponding IP network.
func ParseCIDR(network string) (*net.IPNet, error) {
	_, ip, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	return ip, err
}

// GetAddressCount returns the total number of addresses in the given IP network.
// It accounts for both IPv4 and IPv6 networks, and handles specific cases for certain prefix lengths.
func GetAddressCount(network *net.IPNet) *big.Int {
	prefixLen, bits := network.Mask.Size()

	// Handle specific cases for IPv4 prefix lengths.
	if network.IP.To4() != nil {
		switch prefixLen {
		case 32:
			// A /32 prefix contains a single address.
			return big.NewInt(1)
		case 31:
			// A /31 prefix is used for point-to-point links and contains two addresses.
			return big.NewInt(2)
		}
	}

	// Calculate the total number of addresses based on the prefix length.
	return new(big.Int).Lsh(big.NewInt(1), uint(bits-prefixLen))
}

// GetHostAddressCount returns the number of distinct host addresses in the given IP network.
// It considers the network type (IPv4 or IPv6) and handles edge cases for specific prefix lengths.
// The result excludes the network address and the broadcast address, if applicable.
func GetHostAddressCount(network *net.IPNet) *big.Int {
	prefixLen, bits := network.Mask.Size()

	// Handle edge cases for specific IPv4 prefix lengths.
	if network.IP.To4() != nil {
		switch prefixLen {
		case 32:
			// Single IP address for /32 (e.g., point-to-point link).
			return big.NewInt(1)
		case 31:
			// Two IP addresses for /31 (point-to-point link).
			return big.NewInt(2)
		}
	}

	// Calculate the total number of addresses and subtract 2 (network and broadcast addresses).
	totalAddresses := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits-prefixLen)), nil)
	return totalAddresses.Sub(totalAddresses, big.NewInt(2))
}

// ContainsAddress checks if the given IP network contains the specified IP address.
// It returns true if the address is within the network, otherwise false.
func ContainsAddress(network *net.IPNet, ip net.IP) bool {
	return network.Contains(ip)
}

// Overlaps checks if there is an overlap between two IP networks.
// It returns true if there is any overlap, otherwise false.
func Overlaps(network1, network2 *net.IPNet) bool {
	return network1.Contains(network2.IP) || network2.Contains(network1.IP)
}

// GetNetmask retrieves the netmask associated with the provided IP network.
func GetNetmask(network *net.IPNet) net.IPMask {
	return network.Mask
}

// NetMaskToIPAddress converts a netmask represented as a sequence of bytes
// to its corresponding IP address representation.
func NetMaskToIPAddress(netmask net.IPMask) net.IP {
	return net.IP(netmask)
}

// GetPrefixLength returns the prefix length from the given netmask.
func GetPrefixLength(netmask net.IP) int {
	ones, _ := net.IPMask(netmask).Size()
	return ones
}

// GetBaseAddress returns the base address of the given IP network.
func GetBaseAddress(network *net.IPNet) net.IP {
	return network.IP
}

// GetFirstUsableIPAddress returns the first usable IP address in the given IP network.
func GetFirstUsableIPAddress(network *net.IPNet) (net.IP, error) {
	// If it's an IPv6 network
	if network.IP.To4() == nil {
		ones, bits := network.Mask.Size()
		if ones == bits {
			return nil, errors.New(IPv6NetworkHasNoFirstUsableAddressError)
		}

		// The first address is the first usable address
		firstIP := make(net.IP, len(network.IP))
		copy(firstIP, network.IP)

		return firstIP, nil
	}

	// If it's an IPv4 network, first handle edge cases
	switch ones, _ := network.Mask.Size(); ones {
	case 32:
		return nil, errors.New(IPv4NetworkHasNoFirstUsableAddressError)
	case 31:
		// For /31 network, the current address is the only usable address
		firstIP := make(net.IP, len(network.IP))
		copy(firstIP, network.IP)
		return firstIP, nil
	default:
		// Add 1 to the network address to get the first usable address
		ip := make(net.IP, len(network.IP))
		copy(ip, network.IP)
		ip[3]++ // Add 1 to the last octet

		return ip, nil
	}
}

// GetLastUsableIPAddress returns the last usable IP address in the given IP network.
func GetLastUsableIPAddress(network *net.IPNet) (net.IP, error) {
	// If it's an IPv6 network
	if network.IP.To4() == nil {
		ones, bits := network.Mask.Size()
		if ones == bits {
			return nil, errors.New(IPv6NetworkHasNoLastUsableAddressError)
		}

		// The last address is the last usable address
		lastIP := make(net.IP, len(network.IP))
		copy(lastIP, network.IP)
		for i := range lastIP {
			lastIP[i] |= ^network.Mask[i]
		}

		return lastIP, nil
	}

	// If it's an IPv4 network, first handle edge cases
	switch ones, _ := network.Mask.Size(); ones {
	case 32:
		return nil, errors.New(IPv4NetworkHasNoLastUsableAddressError)
	case 31:
		// For /31 network, the other address is the last usable address
		lastIP := make(net.IP, len(network.IP))
		copy(lastIP, network.IP)
		lastIP[3] |= 1 // Flip the last bit to get the other address
		return lastIP, nil
	default:
		// Subtract 1 from the broadcast address to get the last usable address
		ip := make(net.IP, len(network.IP))
		for i := range ip {
			ip[i] = network.IP[i] | ^network.Mask[i]
		}
		ip[3]-- // Subtract 1 from the last octet

		return ip, nil
	}
}

// GetBroadcastAddress returns the broadcast address of the given IPv4 network, or an error if the IP network is IPv6.
func GetBroadcastAddress(network *net.IPNet) (net.IP, error) {
	if network.IP.To4() == nil {
		// IPv6 networks do not have broadcast addresses.
		return nil, errors.New(IPv6HasNoBroadcastAddressError)
	}

	// Handle edge case for /31 and /32 networks as they have no broadcast address.
	if prefixLen, _ := network.Mask.Size(); helper.ContainsInt([]int{31, 32}, prefixLen) {
		return nil, errors.New(IPv4HasNoBroadcastAddressError)
	}

	ip := make(net.IP, len(network.IP))
	for i := range ip {
		ip[i] = network.IP[i] | ^network.Mask[i]
	}

	return ip, nil
}

// GetMaskWithDivisor calculates the subnet mask for the given divisor and address count.
func GetMaskWithDivisor(divisor int64, addressCount *big.Int, IPv4 bool) (net.IPMask, error) {
	div := big.NewInt(divisor)
	if addressCount.Cmp(div) == -1 || div.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("cannot divide %d addresses into %d divisions", addressCount, div)
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
			return nil, fmt.Errorf("address Space is insufficient for %d subnets", div)
		}
	}
	if subnetPrefix > 126 {
		return nil, fmt.Errorf("address Space is insufficient for %d subnets", div)
	}
	return net.CIDRMask(bits-subnetPrefix, bits), nil
}

// DivideCIDR divides the given IP network into the specified number of subnets.
func DivideCIDR(network *net.IPNet, divisor int64) ([]net.IPNet, error) {
	isIPv4 := helper.IsIPv4Network(network)

	addressCount := GetAddressCount(network)
	newSubnetMask, err := GetMaskWithDivisor(divisor, addressCount, isIPv4)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	networks := make([]net.IPNet, divisor)
	nextAddress := new(net.IPNet)
	nextAddress.IP = network.IP
	nextAddress.Mask = newSubnetMask
	subnetSize := GetAddressCount(nextAddress)
	for i := int64(0); i < divisor; i++ {
		networks[i] = *nextAddress
		ipAsInt := new(big.Int).SetBytes(nextAddress.IP)
		nextAddress.IP = new(big.Int).Add(ipAsInt, subnetSize).Bytes()
	}
	return networks, nil
}
