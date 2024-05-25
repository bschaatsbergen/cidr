# cidr
[![Release](https://github.com/bschaatsbergen/cidr/actions/workflows/goreleaser.yaml/badge.svg)](https://github.com/bschaatsbergen/cidr/actions/workflows/goreleaser.yaml) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/bschaatsbergen/cidr) ![GitHub commits since latest release (by SemVer)](https://img.shields.io/github/commits-since/bschaatsbergen/cidr/latest) [![Go Reference](https://pkg.go.dev/badge/github.com/bschaatsbergen/cidr.svg)](https://pkg.go.dev/github.com/bschaatsbergen/cidr)

CLI to perform various actions on CIDR ranges

## Brew
To install cidr using brew, simply run:

```sh
brew install cidr
```

## Binaries
You can download the [latest binary](https://github.com/bschaatsbergen/cidr/releases/latest) for Linux, MacOS, and Windows.


## Examples

Using `cidr` is very simple.

### Explain a CIDR range

To get more information on a CIDR range:

```
$ cidr explain 10.0.0.0/16
Base Address:            10.0.0.0
Usable Address Range:    10.0.0.1 to 10.0.255.254
Broadcast Address:       10.0.255.255
Address Count:           65,534
Netmask:                 255.255.0.0 (/16 bits)
```

This also works with IPv6 CIDR ranges, for example:

```
$ cidr explain 2001:db8:1234:1a00::/64
Base Address:            2001:db8:1234:1a00::
Usable Address Range:    2001:db8:1234:1a00:: to 2001:db8:1234:1a00:ffff:ffff:ffff:ffff
Address Count:           18,446,744,073,709,551,614
Netmask:                 ffff:ffff:ffff:ffff:: (/64 bits)
```

### Check whether an address belongs to a CIDR range

To check if a CIDR range contains an IP:

```
$ cidr contains 10.0.0.0/16 10.0.14.5
true
```

This also works with IPv6 addresses, for example:

```
$ cidr contains 2001:db8:1234:1a00::/106 2001:db8:1234:1a00::
true
```

### Count distinct host addresses

To get all distinct host addresses part of a given CIDR range:

```
$ cidr count 10.0.0.0/16
65534
```

This also works with a IPv6 CIDR range, for example:

```
$ cidr count 2001:db8:1234:1a00::/106
4194302
```

Or with a large prefix like a point-to-point link CIDR range:

```
$ cidr count 172.16.18.0/31
2
```

### CIDR range intersection

To check if a CIDR range overlaps with another CIDR range:

```
$ cidr overlaps 10.0.0.0/16 10.0.14.0/22
true
```

This also works with IPv6 CIDR ranges, for example:

```
$ cidr overlaps 2001:db8:1111:2222:1::/80 2001:db8:1111:2222:1:1::/96

```

### CIDR division

To divide a cidr range into distinct N distinct networks
## IPV4
```
$ cidr divide 10.0.0.0/16 9
  [Networks]
10.0.0.0/20
10.0.16.0/20
10.0.32.0/20
10.0.48.0/20
10.0.64.0/20
10.0.80.0/20
10.0.96.0/20
10.0.112.0/20
10.0.128.0/20
```

## IPV6
```
$ cidr divide 2001:db8:1111:2222:1::/80 9
  [Networks]
2001:db8:1111:2222:1::/84
2001:db8:1111:2222:1:1000::/84
2001:db8:1111:2222:1:2000::/84
2001:db8:1111:2222:1:3000::/84
2001:db8:1111:2222:1:4000::/84
2001:db8:1111:2222:1:5000::/84
2001:db8:1111:2222:1:6000::/84
2001:db8:1111:2222:1:7000::/84
2001:db8:1111:2222:1:8000::/84

```

You can also use the -u flag to divide the network by desired users/hosts on your network. It assumes a Broadcast and Gateway address in the calculation. So you only have to think of hosts.
The command below illistrate cutting a network into a minimum of 32 hosts, 30 hosts and 12 hosts. And gives you the total possible users/hosts available for your subnet
```
$ cidr d 192.168.0.0/24 -u 32,30,12
  [Networks]            [Used]  [Total]
192.168.0.0/26            32      62
192.168.0.64/27           30      30
192.168.0.96/28           12      14

```
## Contributing

Contributions are highly appreciated and always welcome.
Have a look through existing [Issues](https://github.com/bschaatsbergen/cidr/issues) and [Pull Requests](https://github.com/bschaatsbergen/cidr/pulls) that you could help with.
