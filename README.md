# cidr
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/bschaatsbergen/cidr) [![Go Reference](https://pkg.go.dev/badge/github.com/bschaatsbergen/cidr.svg)](https://pkg.go.dev/github.com/bschaatsbergen/cidr)

Simplifies IPv4/IPv6 CIDR network prefix management with counting, overlap checking, explanation, and subdivision

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
Usable Address Range:    10.0.0.1 to 10.0.255.254 (65,534)
Broadcast Address:       10.0.255.255
Addresses:               65,536
Netmask:                 255.255.0.0 (/16 bits)
```

This also works with IPv6 CIDR ranges, for example:

```
$ cidr explain 2001:db8:1234:1a00::/110
Base Address:            2001:db8:1234:1a00::
Usable Address Range:    2001:db8:1234:1a00:: to 2001:db8:1234:1a00::3:ffff (262,142)
Addresses:               262,144
Netmask:                 ffff:ffff:ffff:ffff:ffff:ffff:fffc:0 (/110 bits)
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

### Count

To get a count of all addresses in a CIDR range:

```
$ cidr count 10.0.0.0/16
65536
```

This also works with a IPv6 CIDR range, for example:

```
$ cidr count 2001:db8:1234:1a00::/106
4194304
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
true
```

### CIDR division

To divide a CIDR range into N distinct networks:

```
$ cidr divide 10.0.0.0/16 9
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

This also works with IPv6 CIDR ranges, for example:

```
$ cidr divide 2001:db8:1111:2222:1::/80 9
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

## Contributing

Contributions are highly appreciated and always welcome.
Have a look through existing [Issues](https://github.com/bschaatsbergen/cidr/issues) and [Pull Requests](https://github.com/bschaatsbergen/cidr/pulls) that you could help with.
