// Copyright (c) Bruno Schaatsbergen
// SPDX-License-Identifier: MIT

package helper

import "net"

// isIPv4Network checks if the given network is an IPv4 network.
// It returns true if the network is an IPv4 network, otherwise false.
func IsIPv4Network(network *net.IPNet) bool {
	return network.IP.To4() != nil
}

// isIPv6Network checks if the given network is an IPv6 network.
// It returns true if the network is an IPv6 network, otherwise false.
func IsIPv6Network(network *net.IPNet) bool {
	return network.IP.To16() != nil
}
