package core

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressCount(t *testing.T) {
	IPv4CIDR, err := ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := ParseCIDR("2001:db8:1234:1a00::/106")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	largeIPv4PrefixCIDR, err := ParseCIDR("172.16.18.0/31")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name          string
		cidr          *net.IPNet
		expectedCount uint64
	}{
		{
			name:          "Return the count of all distinct host addresses in a common IPv4 CIDR",
			cidr:          IPv4CIDR,
			expectedCount: 65534,
		},
		{
			name:          "Return the count of all distinct host addresses in a common IPv6 CIDR",
			cidr:          IPv6CIDR,
			expectedCount: 4194302,
		},
		{
			name:          "Return the count of all distinct host addresses in an uncommon (large prefix) IPv4 CIDR",
			cidr:          largeIPv4PrefixCIDR,
			expectedCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := AddressCount(tt.cidr)
			assert.Equal(t, int(tt.expectedCount), int(count), "Both address counts should be equal")
		})
	}
}

func TestOverlaps(t *testing.T) {
	firstIPv4CIDR, err := ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	secondIPv4CIDR, err := ParseCIDR("10.0.14.0/22")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	thirdIPv4CIDR, err := ParseCIDR("10.1.0.0/28")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	firstIPv6CIDR, err := ParseCIDR("2001:db8:1111:2222:1::/80")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	secondIPv6CIDR, err := ParseCIDR("2001:db8:1111:2222:1:1::/96")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name     string
		cidrA    *net.IPNet
		cidrB    *net.IPNet
		overlaps bool
	}{
		{
			name:     "2 IPv4 CIDR ranges should overlap",
			cidrA:    firstIPv4CIDR,
			cidrB:    secondIPv4CIDR,
			overlaps: true,
		},
		{
			name:     "2 IPv4 CIDR ranges should NOT overlap",
			cidrA:    firstIPv4CIDR,
			cidrB:    thirdIPv4CIDR,
			overlaps: false,
		},
		{
			name:     "2 IPv6 CIDR ranges should overlap",
			cidrA:    firstIPv6CIDR,
			cidrB:    secondIPv6CIDR,
			overlaps: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overlaps := Overlaps(tt.cidrA, tt.cidrB)
			assert.Equal(t, tt.overlaps, overlaps, "Foobar")
		})
	}
}
