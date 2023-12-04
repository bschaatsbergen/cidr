package core

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAddressCount(t *testing.T) {
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

	largestIPv4PrefixCIDR, err := ParseCIDR("172.16.18.0/32")
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
		{
			name:          "Return the count of all distinct host addresses in an uncommon (largest prefix) IPv4 CIDR",
			cidr:          largestIPv4PrefixCIDR,
			expectedCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := GetAddressCount(tt.cidr)
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
			assert.Equal(t, tt.overlaps, overlaps, "Given CIDRs should overlap")
		})
	}
}

func TestContainsAddress(t *testing.T) {
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

	tests := []struct {
		name     string
		cidr     *net.IPNet
		ip       net.IP
		contains bool
	}{
		{
			name:     "IPv4 CIDR that does contain an IPv4 IP",
			cidr:     IPv4CIDR,
			ip:       net.ParseIP("10.0.14.5"),
			contains: true,
		},
		{
			name:     "IPv4 CIDR that does NOT contain an IPv4 IP",
			cidr:     IPv4CIDR,
			ip:       net.ParseIP("10.1.55.5"),
			contains: false,
		},
		{
			name:     "IPv6 CIDR that does contain an IPv6 IP",
			cidr:     IPv6CIDR,
			ip:       net.ParseIP("2001:db8:1234:1a00::"),
			contains: true,
		},
		{
			name:     "IPv6 CIDR that does NOT contain an IPv6 IP",
			cidr:     IPv6CIDR,
			ip:       net.ParseIP("2001:af1:1222:1a20::"),
			contains: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overlaps := ContainsAddress(tt.cidr, tt.ip)
			assert.Equal(t, tt.contains, overlaps, "Given IP address should be part of the given CIDR")
		})
	}
}

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		name    string
		cidrStr string
		wantErr bool
	}{
		{
			name:    "Parse a valid IPv4 CIDR",
			cidrStr: "10.0.0.0/16",
			wantErr: false,
		},
		{
			name:    "Parse a valid IPv6 CIDR",
			cidrStr: "2001:db8:1234:1a00::/106",
			wantErr: false,
		},
		{
			name:    "Parse an invalid IPv4 CIDR",
			cidrStr: "356.356.356.356/16",
			wantErr: true,
		},
		{
			name:    "Parse an invalid IPv6 CIDR",
			cidrStr: "2001:db8:1234:1a00::/129",
			wantErr: true,
		},
		{
			name:    "Parse an empty CIDR",
			cidrStr: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCIDR(tt.cidrStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCIDR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
