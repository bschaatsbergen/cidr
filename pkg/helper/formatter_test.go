// Copyright (c) Bruno Schaatsbergen
// SPDX-License-Identifier: MIT

package helper_test

import (
	"testing"

	"github.com/bschaatsbergen/cidr/pkg/helper"
)

func TestFormatNumber(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"1234", "1,234"},
		{"12345", "12,345"},
		{"123456", "123,456"},
		{"1234567", "1,234,567"},
		{"12345678", "12,345,678"},
		{"123456789", "123,456,789"},
		{"1234567890", "1,234,567,890"},
		{"", "<nil>"},      // Empty input
		{"abc", "<nil>"},   // Non-numeric input
		{"1.23", "<nil>"},  // Decimal input
		{"-1234", "<nil>"}, // Negative input
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := helper.FormatNumber(tc.input)
			if result != tc.expected {
				t.Errorf("For input '%s', expected '%s' but got '%s'", tc.input, tc.expected, result)
			}
		})
	}
}
