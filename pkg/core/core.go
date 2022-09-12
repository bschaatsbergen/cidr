package core

import (
	"net"
)

func AddressCount(network *net.IPNet) uint64 {
	prefixLen, bits := network.Mask.Size()

	// Check if network is IPv4 or IPv6
	if network.Mask != nil {
		// Handle edge cases
		switch prefixLen {
			case 32: return 1
			case 31: return 2
		}
	}

	// Remember to subtract the network address and broadcast address
	return 1 << (uint64(bits) - uint64(prefixLen)) - 2
}

func ParseCIDR(network string) (*net.IPNet, error) {
	_, ip, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	return ip, err
}

func ParseIP(ip string) net.IP {
	return net.ParseIP(ip)
}

func ContainsAddress(network *net.IPNet, ip net.IP) bool {
	return network.Contains(ip)
}

func Overlaps(network1, network2 *net.IPNet) bool {
	return network1.Contains(network2.IP) || network2.Contains(network1.IP)
}
