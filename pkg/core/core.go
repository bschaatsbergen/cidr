package core

import (
	"net"
)

func AddressCount(network *net.IPNet) uint64 {
	prefixLen, bits := network.Mask.Size()
	return 1 << (uint64(bits) - uint64(prefixLen))
}

func ParseCIDR(network string) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	return ipNet, err
}
