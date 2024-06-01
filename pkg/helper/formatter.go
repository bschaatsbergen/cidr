// Copyright (c) Bruno Schaatsbergen
// SPDX-License-Identifier: MIT

package helper


// Takes a numbered string [123456789].
// Outputs -> [123,456,789].
func FormatNumber(s string) string {
	length := len(s)
	if length == 0 { return "<nil>"}
	newLength := length + (length - 1) / 3
	newNumber := make([]rune, newLength)
	commas := 0
	for i, c := range s {
		if c < '0' || c > '9' {return "<nil>"}
		if (length - i) % 3 == 0 && i != 0 {
			newNumber[i+commas] = ','
			commas += 1
		}
		newNumber[i+commas] = c

	}
	return string(newNumber)
}
