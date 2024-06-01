// Copyright (c) Bruno Schaatsbergen
// SPDX-License-Identifier: MIT

package core

import (
	"errors"
	"math/big"
	"net"
	"fmt"
	"math"

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

// GetAddressCount returns the number of addresses in the given IP network.
// It considers the network type (IPv4 or IPv6) and handles edge cases for specific prefix lengths.
// The result excludes the network address and broadcast address.
func GetAddressCount(network *net.IPNet) *big.Int {
	prefixLen, bits := network.Mask.Size()

	// Handle edge cases for specific IPv4 prefix lengths.
	if network.Mask != nil && network.IP.To4() != nil {
		switch prefixLen {
		case 32:
			return big.NewInt(1)
		case 31:
			return big.NewInt(2)
		}
	}

	return big.NewInt(0).Lsh(big.NewInt(1), uint(bits-prefixLen))
}

// GetNextAddress retrieves the next available IPNet with your desired CIDR.
func GetNextAddress(network *net.IPNet, cidr net.IPMask) (net.IPNet, error) {
	addressCount := GetAddressCount(network)
	var currentIP *big.Int
	if helper.IsIPv4Network(network) {
		currentIP = big.NewInt(0).SetBytes(network.IP.To4())
	} else {
		currentIP = big.NewInt(0).SetBytes(network.IP.To16())
	}
	nextAddressNum := new(big.Int).Add(addressCount, currentIP)

	nextAddress := net.IP(nextAddressNum.Bytes())
	return net.IPNet{
		IP:   nextAddress,
		Mask: cidr,
	}, nil

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


// Returns the net.IPMask necessary for the provided divisor.
// Errors if address space if insufficient or division is not possible.
func GetMaskWithDivisor(divisor int64, addressCount *big.Int, IPv4 bool) (net.IPMask, error) {
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


// Divides the given network into N smaller networks.
// Errors if division is not possible.
func DivideCidr(network *net.IPNet, divisor int64) ([]net.IPNet, error) {
	isIPv4 := helper.IsIPv4Network(network)

	addressCount := GetAddressCount(network)
	cidrWack, err := GetMaskWithDivisor(divisor, addressCount, isIPv4)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	networks := make([]net.IPNet, divisor)
	nextAddress := new(net.IPNet)
	nextAddress.IP = network.IP
	nextAddress.Mask = cidrWack
	subnetSize := GetAddressCount(nextAddress)
	for i := int64(0); i < divisor; i++ {
		networks[i] = *nextAddress
		ipAsInt := new(big.Int).SetBytes(nextAddress.IP)
		nextAddress.IP = new(big.Int).Add(ipAsInt, subnetSize).Bytes()
	}
	return networks, nil
}

// Given an IP Network and a list of desired hosts, splits the network into a configuration that accommodates the user space.
// Eg. Splitting a 192.168.0.0/24 where I need a network with 20 hosts, 10 hosts, and 30users.
// Gets the /27 /28 and /27 combination.
func DivideCidrHosts(network *net.IPNet, desiredHosts []int64) ([]net.IPNet, error) {
	var bits int
	if helper.IsIPv4Network(network) {
		bits = net.IPv4len * 8
	} else {
		bits = net.IPv6len * 8
	}
	networks := make([]net.IPNet, len(desiredHosts))
	nextIP := network.IP
	for i, hosts := range desiredHosts {
		subnetPrefix := net.CIDRMask(bits-int(hosts), bits)
		nextNetwork := net.IPNet{
			IP:   nextIP,
			Mask: subnetPrefix,
		}
		networks[i] = nextNetwork
		nextNetwork, err := GetNextAddress(&nextNetwork, subnetPrefix)
		if err != nil {
			return nil, err
		}
		nextIP = nextNetwork.IP
	}
	return networks, nil
}


func ValidateHostSpace(network *net.IPNet, desiredHostsPerSubnet []int64) ([]int64, error) {
	requiredAddressSpace := int64(0)
	totalAddressSpace := GetAddressCount(network)
	desiredHosts := make([]int64, len(desiredHostsPerSubnet))
	for i, hosts := range desiredHostsPerSubnet {
		subnetExponent := math.Ceil(math.Log2(float64(hosts) + 2.0))
		addressSpace := int64(math.Pow(2, subnetExponent))
		requiredAddressSpace += addressSpace
		desiredHosts[i] = int64(subnetExponent)
	}
	requiredAddressSpaceBig := big.NewInt(requiredAddressSpace)
	if totalAddressSpace.Cmp(requiredAddressSpaceBig) < 0 {
		return nil, fmt.Errorf("Total address space is: %s but desired Hosts requires %d addresses\n", totalAddressSpace.String(), requiredAddressSpace)
	}
	return desiredHosts, nil
}
