package core

import (
	"net"
)

func AddressCount(networkAdress *net.IPNet) uint64 {
	prefixLen, bits := networkAdress.Mask.Size()
	return 1 << (uint64(bits) - uint64(prefixLen))
}

func ParseCIDR(networkAdress string) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(networkAdress)
	if err != nil {
		return nil, err
	}
	return ipNet, err
}

func ParseIP(ipAddress string) net.IP {
	return net.ParseIP(ipAddress)
}

func ContainsAddress(networkAdress *net.IPNet, ipAddress net.IP) bool {
	return networkAdress.Contains(ipAddress)
}
