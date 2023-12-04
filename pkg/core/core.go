package core

import (
	"net"
)

// GetAddressCount returns the number of usable addresses in the given IP network.
// It considers the network type (IPv4 or IPv6) and handles edge cases for specific prefix lengths.
// The result excludes the network address and broadcast address.
func GetAddressCount(network *net.IPNet) uint64 {
	prefixLen, bits := network.Mask.Size()

	// Handle edge cases for specific IPv4 prefix lengths.
	if network.Mask != nil && network.IP.To4() != nil {
		switch prefixLen {
		case 32:
			return 1
		case 31:
			return 2
		}
	}

	// Subtract the network address and broadcast address (2) from the total number of addresses.
	return 1<<(uint64(bits)-uint64(prefixLen)) - 2
}

// ParseCIDR parses the given CIDR notation string and returns the corresponding IP network.
func ParseCIDR(network string) (*net.IPNet, error) {
	_, ip, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	return ip, err
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

func ListCIDR(network string) ([]string, error) {
	var hosts []string
	ip, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); increment(ip) {
		hosts = append(hosts, ip.String())
	}
	return hosts, nil
}

func increment(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}
