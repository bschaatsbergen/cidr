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
		case 32:
			return 1
		case 31:
			return 2
		}
	}

	// Remember to subtract the network address and broadcast address
	return 1<<(uint64(bits)-uint64(prefixLen)) - 2
}

func ParseCIDR(network string) (*net.IPNet, error) {
	_, ip, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	return ip, err
}

func ContainsAddress(network *net.IPNet, ip net.IP) bool {
	return network.Contains(ip)
}

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
