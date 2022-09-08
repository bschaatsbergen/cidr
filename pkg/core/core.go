package core

import (
	"net"
)

func AddressCount(n *net.IPNet) uint64 {
	prefixLen, bits := n.Mask.Size()
	return 1 << (uint64(bits) - uint64(prefixLen))
}

func ParseCIDR(n string) (*net.IPNet, error) {
	_, ip, err := net.ParseCIDR(n)
	if err != nil {
		return nil, err
	}
	return ip, err
}

func ParseIP(i string) net.IP {
	return net.ParseIP(i)
}

func ContainsAddress(n *net.IPNet, i net.IP) bool {
	return n.Contains(i)
}

func Overlaps(n1, n2 *net.IPNet) bool {
	return n1.Contains(n2.IP) || n2.Contains(n1.IP)
}
