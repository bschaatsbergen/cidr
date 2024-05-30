// Copyright (c) Bruno Schaatsbergen
// SPDX-License-Identifier: MIT

package cmd_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/bschaatsbergen/cidr/cmd"
)

func TestDivideCidr(t *testing.T) {
	testCases := []struct {
		cidr      string
		divisor   int64
		expected  []string // Expected CIDRs in string format
		shouldErr bool
	}{
		{"10.0.0.0/16", 2, []string{"10.0.0.0/17", "10.0.128.0/17"}, false},
		{"192.168.0.0/24", 4, []string{"192.168.0.0/26", "192.168.0.64/26", "192.168.0.128/26", "192.168.0.192/26"}, false},
		{"2001:db8::/32", 3, []string{"2001:db8::/34", "2001:db8:4000::/34", "2001:db8:8000::/34"}, false},
		// Error cases
		{"10.0.0.0/16", 0, nil, true},    // Divisor is zero
		{"2001:db8::/128", 3, nil, true}, // Cannot divide /128
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

func TestDivideCidrHosts(t *testing.T) {
	testCases := []struct {
		inputCIDR    string
		desiredUsers []int64
		expectedNets []string
		expectedErr  error
	}{
		{"192.168.0.0/24", []int64{20, 10, 30}, []string{"192.168.0.0/27", "192.168.0.32/28", "192.168.0.48/27"}, nil},
		{"10.0.0.0/16", []int64{1000, 500}, []string{"10.0.0.0/22", "10.0.4.0/23"}, nil},
		{"2001:db8::/32", []int64{50000}, []string{"2001:db8::/112"}, nil},
		{"192.168.0.0/24", []int64{}, []string{}, nil},                                                                              // Edge case: empty desiredUsers
		{"192.168.0.0/24", []int64{257}, nil, fmt.Errorf("Total address space is: 256 but desired Users requires 512 addresses\n")}, // Edge case: not enough addresses
	}

	for _, tc := range testCases {
		_, ipNet, err := net.ParseCIDR(tc.inputCIDR)
		if err != nil && tc.expectedErr == nil {
			t.Fatalf("Failed to parse CIDR: %v", err)
		}

		validatedUsers, err := cmd.ValidateUserSpace(ipNet, tc.desiredUsers)
		if tc.expectedErr != nil {
			if err == nil || err.Error() != tc.expectedErr.Error() {
				t.Errorf("A: For %s, users %v: Expected error '%v', got '%v'", tc.inputCIDR, tc.desiredUsers, tc.expectedErr, err)
			}
			continue
		}

		networks, err := cmd.DivideCidrHosts(ipNet, validatedUsers)

		if tc.expectedErr != nil {
			if err == nil || err.Error() != tc.expectedErr.Error() {
				t.Errorf("A: For %s, users %v: Expected error '%v', got '%v'", tc.inputCIDR, tc.desiredUsers, tc.expectedErr, err)
			}
			continue
		}

		if err != nil {
			t.Errorf("B: For %s, users %v: Unexpected error '%v'", tc.inputCIDR, tc.desiredUsers, err)
			continue
		}

		netStrings := make([]string, len(networks))
		for i, n := range networks {
			netStrings[i] = n.String()
		}

		if !equalStringSlices(netStrings, tc.expectedNets) {
			fmt.Println(networks)
			t.Errorf("C: For %s, users %v: Expected subnets '%v', got '%v'", tc.inputCIDR, tc.desiredUsers, tc.expectedNets, netStrings)
		}
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
