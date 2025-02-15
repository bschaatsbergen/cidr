package core_test

import (
	"math/big"
	"net"
	"testing"

	"github.com/bschaatsbergen/cidr/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestGetAddressCount(t *testing.T) {
	IPv4CIDR, err := core.ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := core.ParseCIDR("2001:db8:1234:1a00::/106")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	largeIPv4PrefixCIDR, err := core.ParseCIDR("172.16.18.0/31")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	largestIPv4PrefixCIDR, err := core.ParseCIDR("172.16.18.0/32")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name          string
		cidr          *net.IPNet
		expectedCount *big.Int
	}{
		{
			name:          "Return the count of all addresses in a common IPv4 CIDR",
			cidr:          IPv4CIDR,
			expectedCount: big.NewInt(65536),
		},
		{
			name:          "Return the count of all addresses in a common IPv6 CIDR",
			cidr:          IPv6CIDR,
			expectedCount: big.NewInt(4194304),
		},
		{
			name:          "Return the count of all addresses in an uncommon (large prefix) IPv4 CIDR",
			cidr:          largeIPv4PrefixCIDR,
			expectedCount: big.NewInt(2),
		},
		{
			name:          "Return the count of all addresses in an uncommon (largest prefix) IPv4 CIDR",
			cidr:          largestIPv4PrefixCIDR,
			expectedCount: big.NewInt(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := core.GetAddressCount(tt.cidr)
			assert.Equal(t, tt.expectedCount, count, "Both address counts should be equal")
		})
	}
}

func TestGetHostAddressCount(t *testing.T) {
	IPv4CIDR, err := core.ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := core.ParseCIDR("2001:db8:1234:1a00::/106")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	largeIPv4PrefixCIDR, err := core.ParseCIDR("172.16.18.0/31")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	largestIPv4PrefixCIDR, err := core.ParseCIDR("172.16.18.0/32")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name          string
		cidr          *net.IPNet
		expectedCount *big.Int
	}{
		{
			name:          "Return the count of all distinct host addresses in a common IPv4 CIDR",
			cidr:          IPv4CIDR,
			expectedCount: big.NewInt(65534),
		},
		{
			name:          "Return the count of all distinct host addresses in a common IPv6 CIDR",
			cidr:          IPv6CIDR,
			expectedCount: big.NewInt(4194302),
		},
		{
			name:          "Return the count of all distinct host addresses in an uncommon (large prefix) IPv4 CIDR",
			cidr:          largeIPv4PrefixCIDR,
			expectedCount: big.NewInt(2),
		},
		{
			name:          "Return the count of all distinct host addresses in an uncommon (largest prefix) IPv4 CIDR",
			cidr:          largestIPv4PrefixCIDR,
			expectedCount: big.NewInt(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := core.GetHostAddressCount(tt.cidr)
			assert.Equal(t, tt.expectedCount, count, "Both address counts should be equal")
		})
	}
}

func TestOverlaps(t *testing.T) {
	firstIPv4CIDR, err := core.ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	secondIPv4CIDR, err := core.ParseCIDR("10.0.14.0/22")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	thirdIPv4CIDR, err := core.ParseCIDR("10.1.0.0/28")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	firstIPv6CIDR, err := core.ParseCIDR("2001:db8:1111:2222:1::/80")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	secondIPv6CIDR, err := core.ParseCIDR("2001:db8:1111:2222:1:1::/96")
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
			overlaps := core.Overlaps(tt.cidrA, tt.cidrB)
			assert.Equal(t, tt.overlaps, overlaps, "Given CIDRs should overlap")
		})
	}
}

func TestContainsAddress(t *testing.T) {
	IPv4CIDR, err := core.ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := core.ParseCIDR("2001:db8:1234:1a00::/106")
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
			overlaps := core.ContainsAddress(tt.cidr, tt.ip)
			assert.Equal(t, tt.contains, overlaps, "Given IP address should be part of the given CIDR")
		})
	}
}

func TestGetPrefixLength(t *testing.T) {
	IPv4CIDR, err := core.ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := core.ParseCIDR("2001:db8:1234:1a00::/106")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name                 string
		netmask              net.IP
		expectedPrefixLength int
	}{
		{
			name:                 "Get the prefix length of an IPv4 netmask",
			netmask:              net.IP(IPv4CIDR.Mask),
			expectedPrefixLength: 16,
		},
		{
			name:                 "Get the prefix length of an IPv6 netmask",
			netmask:              net.IP(IPv6CIDR.Mask),
			expectedPrefixLength: 106,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefixLength := core.GetPrefixLength(tt.netmask)
			assert.Equal(t, tt.expectedPrefixLength, prefixLength, "Prefix length is not correct")
		})
	}
}

func TestGetNetMask(t *testing.T) {
	IPv4CIDR, err := core.ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := core.ParseCIDR("2001:db8:1234:1a00::/106")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name            string
		CIDR            *net.IPNet
		expectedNetmask net.IPMask
	}{
		{
			name:            "Get the netmask of an IPv4 CIDR range",
			CIDR:            IPv4CIDR,
			expectedNetmask: net.IPMask{0xff, 0xff, 0x00, 0x00},
		},
		{
			name:            "Get the netmask of an IPv6 CIDR range",
			CIDR:            IPv6CIDR,
			expectedNetmask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xc0, 0x0, 0x0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			netmask := core.GetNetmask(tt.CIDR)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			assert.Equal(t, tt.expectedNetmask, netmask, "Netmask is not correct")
		})
	}
}

func TestGetFirstUsableIPAddress(t *testing.T) {
	IPv4CIDR, err := core.ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := core.ParseCIDR("2001:db8:1234:1a00::/106")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name                         string
		CIDR                         *net.IPNet
		expectedFirstUsableIPAddress net.IP
	}{
		{
			name:                         "Get the first usable IP address of an IPv4 CIDR range",
			CIDR:                         IPv4CIDR,
			expectedFirstUsableIPAddress: net.ParseIP("10.0.0.1").To4(),
		},
		{
			name:                         "Get the first usable IP address of an IPv6 CIDR range",
			CIDR:                         IPv6CIDR,
			expectedFirstUsableIPAddress: net.ParseIP("2001:db8:1234:1a00::").To16(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstUsableIPAddress, err := core.GetFirstUsableIPAddress(tt.CIDR)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			assert.Equal(t, tt.expectedFirstUsableIPAddress, firstUsableIPAddress, "First usable IP address is not correct")
		})
	}
}

func TestGetLastUsableIPAddress(t *testing.T) {
	IPv4CIDR, err := core.ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := core.ParseCIDR("2001:db8:1234:1a00::/106")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name                        string
		CIDR                        *net.IPNet
		expectedLastUsableIPAddress net.IP
	}{
		{
			name:                        "Get the last usable IP address of an IPv4 CIDR range",
			CIDR:                        IPv4CIDR,
			expectedLastUsableIPAddress: net.ParseIP("10.0.255.254").To4(),
		},
		{
			name:                        "Get the last usable IP address of an IPv6 CIDR range",
			CIDR:                        IPv6CIDR,
			expectedLastUsableIPAddress: net.ParseIP("2001:db8:1234:1a00::3f:ffff").To16(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastUsableIPAddress, err := core.GetLastUsableIPAddress(tt.CIDR)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			assert.Equal(t, tt.expectedLastUsableIPAddress, lastUsableIPAddress, "Last usable IP address is not correct")
		})
	}
}

func TestGetBroadcastAddress(t *testing.T) {
	IPv4CIDR, err := core.ParseCIDR("10.0.0.0/16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv4CIDRWithNoBroadcastAddress, err := core.ParseCIDR("10.0.0.0/31")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := core.ParseCIDR("2001:db8:1234:1a00::/106")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name                     string
		CIDR                     *net.IPNet
		expectedBroadcastAddress net.IP
		wantErr                  bool
	}{
		{
			name:                     "Get the broadcast IP address of an IPv4 CIDR range",
			CIDR:                     IPv4CIDR,
			expectedBroadcastAddress: net.ParseIP("10.0.255.255").To4(),
			wantErr:                  false,
		},
		{
			name:                     "Get the broadcast IP address of an IPv4 CIDR range that has no broadcast address",
			CIDR:                     IPv4CIDRWithNoBroadcastAddress,
			expectedBroadcastAddress: nil,
			wantErr:                  true,
		},
		{
			name:                     "Get the broadcast IP address of an IPv6 CIDR range",
			CIDR:                     IPv6CIDR,
			expectedBroadcastAddress: nil,
			wantErr:                  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			broadcastAddress, err := core.GetBroadcastAddress(tt.CIDR)
			if err != nil {
				assert.Equal(t, tt.wantErr, true, "Expected error when getting broadcast address, but got none")
			} else {
				assert.Equal(t, tt.expectedBroadcastAddress, broadcastAddress, "Broadcast IP address is not correct")
			}
		})
	}
}

func TestGetBaseAddress(t *testing.T) {
	IPv4CIDR, err := core.ParseCIDR("192.168.90.4/30")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	IPv6CIDR, err := core.ParseCIDR("4a00:db8:1234:1a00::/127")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	tests := []struct {
		name                string
		cidr                *net.IPNet
		expectedBaseAddress net.IP
	}{
		{
			name:                "Get the base address of an IPv4 CIDR",
			cidr:                IPv4CIDR,
			expectedBaseAddress: IPv4CIDR.IP,
		},
		{
			name:                "Get the base address of an IPv6 CIDR",
			cidr:                IPv6CIDR,
			expectedBaseAddress: IPv6CIDR.IP,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseAddress := core.GetBaseAddress(tt.cidr)
			assert.Equal(t, tt.expectedBaseAddress, baseAddress, "Base address is not correct")
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
			_, err := core.ParseCIDR(tt.cidrStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCIDR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDivideCIDR(t *testing.T) {
	tests := []struct {
		name      string
		cidr      string
		divisor   int64
		expected  []string
		shouldErr bool
	}{
		{
			name:      "Divide IPv4 CIDR into 2 subnets",
			cidr:      "10.0.0.0/16",
			divisor:   2,
			expected:  []string{"10.0.0.0/17", "10.0.128.0/17"},
			shouldErr: false,
		},
		{
			name:      "Divide IPv4 CIDR into 4 subnets",
			cidr:      "192.168.0.0/24",
			divisor:   4,
			expected:  []string{"192.168.0.0/26", "192.168.0.64/26", "192.168.0.128/26", "192.168.0.192/26"},
			shouldErr: false,
		},
		{
			name:      "Divide IPv6 CIDR into 3 subnets",
			cidr:      "2001:db8::/32",
			divisor:   3,
			expected:  []string{"2001:db8::/34", "2001:db8:4000::/34", "2001:db8:8000::/34"},
			shouldErr: false,
		},
		{
			name:      "Error case: Divisor is zero",
			cidr:      "10.0.0.0/16",
			divisor:   0,
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "Error case: Cannot divide /128 CIDR",
			cidr:      "2001:db8::/128",
			divisor:   3,
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ipNet, err := net.ParseCIDR(tt.cidr)
			assert.NoError(t, err, "Unexpected error parsing CIDR: %v", err)

			subnets, err := core.DivideCIDR(ipNet, tt.divisor)
			if tt.shouldErr {
				assert.Error(t, err, "Expected error but got none")
				return
			}

			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.Equal(t, len(tt.expected), len(subnets), "Incorrect number of subnets")

			for i, expectedCIDR := range tt.expected {
				assert.Equal(t, expectedCIDR, subnets[i].String(), "Incorrect subnet at index %d", i)
			}
		})
	}
}
