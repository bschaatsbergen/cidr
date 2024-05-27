package cmd_test

import (
	"math/big"
	"net"
	"testing"

	"github.com/bschaatsbergen/cidr/pkg/core"
)

func TestGetAddressCount(t *testing.T) {
	testCases := []struct {
		cidr     string
		expected *big.Int
	}{
		// Common IPv4 cases
		{"10.0.0.0/24", big.NewInt(256)},
		{"192.168.1.0/28", big.NewInt(16)},
		{"172.16.0.0/12", big.NewInt(1048576)},
		// IPv4 edge cases
		{"192.168.0.0/31", big.NewInt(2)},
		{"10.10.10.10/32", big.NewInt(1)},
		// IPv6 cases
		{"2001:db8::/32", big.NewInt(0).Lsh(big.NewInt(1), 96)}, // [128 - 32 -> 2^96 ]
		{"fe80::/64", big.NewInt(0).Lsh(big.NewInt(1), 64)},     // [128 - 64 -> 2^64 ]
		{"::1/128", big.NewInt(1)},
	}

	for _, tc := range testCases {
		_, ipNet, err := net.ParseCIDR(tc.cidr)
		if err != nil {
			t.Fatalf("Failed to parse CIDR %s: %v", tc.cidr, err)
		}

		actual := core.GetAddressCount(ipNet)
		if actual.Cmp(tc.expected) != 0 {
			t.Errorf("For CIDR %s, expected address count %d, but got %d", tc.cidr, tc.expected, actual)
		}
	}
}
