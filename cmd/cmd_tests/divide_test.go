package cmd_test

import (
	"net"
	"testing"

	"github.com/bschaatsbergen/cidr/cmd"
)

func TestDivideCidr(t *testing.T) {
    testCases := []struct {
        cidr      string
        divisor   uint64
        expected  []string // Expected CIDRs in string format
        shouldErr bool
    }{
        {"10.0.0.0/16", 2, []string{"10.0.0.0/17", "10.0.128.0/17"}, false},
        {"192.168.0.0/24", 4, []string{"192.168.0.0/26", "192.168.0.64/26", "192.168.0.128/26", "192.168.0.192/26"}, false},
        {"2001:db8::/32", 3, []string{"2001:db8::/34", "2001:db8:4000::/34", "2001:db8:8000::/34"}, false},
        // Error cases
        {"10.0.0.0/16", 0, nil, true},      // Divisor is zero
        {"2001:db8::/128", 3, nil, true},   // Cannot divide /128
    }

    for _, tc := range testCases {
        _, ipNet, err := net.ParseCIDR(tc.cidr)
        if err != nil && !tc.shouldErr {
            t.Errorf("Unexpected error parsing CIDR: %s", err)
            continue
        }

        subnets, err := cmd.DivideCidr(ipNet, tc.divisor)


        if tc.shouldErr {
            if err == nil {
                t.Errorf("Expected error for CIDR %s, divisor %d, but got none", tc.cidr, tc.divisor)
            }
            continue
        }

        if err != nil {
            t.Errorf("Unexpected error for CIDR %s, divisor %d: %s", tc.cidr, tc.divisor, err)
            continue
        }

        if len(subnets) != len(tc.expected) {
            t.Errorf("Incorrect number of subnets for CIDR %s, divisor %d: expected %d, got %d", tc.cidr, tc.divisor, len(tc.expected), len(subnets))
            continue
        }

        for i, expectedCIDR := range tc.expected {
            if subnets[i].String() != expectedCIDR {
                t.Errorf("Incorrect subnet %d for CIDR %s, divisor %d: expected %s, got %s", i, tc.cidr, tc.divisor, expectedCIDR, subnets[i].String())
            }
        }
    }
}
