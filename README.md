# cidr
[![Release](https://github.com/bschaatsbergen/cidr/actions/workflows/goreleaser.yaml/badge.svg)](https://github.com/bschaatsbergen/cidr/actions/workflows/goreleaser.yaml) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/bschaatsbergen/cidr) ![GitHub commits since latest release (by SemVer)](https://img.shields.io/github/commits-since/bschaatsbergen/cidr/latest) [![Go Reference](https://pkg.go.dev/badge/github.com/bschaatsbergen/cidr.svg)](https://pkg.go.dev/github.com/bschaatsbergen/cidr)

A CLI to perform various actions on CIDR ranges

## Brew
To install cidr using brew, simply do the below.

```sh
brew tap bschaatsbergen/cidr
brew install cidr
```

## Binaries
You can download the [latest binary](https://github.com/bschaatsbergen/cidr/releases/latest) for Linux, MacOS, and Windows.


## Examples

Using `cidr` is very simple.

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

Or with large prefixes like point-to-point links:

```
$ cidr count count 172.16.18.0/31
2
```

This also works with very small CIDR ranges, like a point-to-point link:

```
$ cidr count 2001:db8:1234:1a00::/106
4194304
```

### CIDR range intersection

To check if a CIDR range overlaps with another CIDR range:

```
$ cidr overlaps 10.0.0.0/16 10.0.14.0/22
true
```

This also works with IPv6 addresses, for example:

```
$ cidr overlaps 2001:db8:1111:2222:1::/80 2001:db8:1111:2222:1:1::/96
true
```

## Contributing

Contributions are highly appreciated and always welcome.
Have a look through existing [Issues](https://github.com/bschaatsbergen/cidr/issues) and [Pull Requests](https://github.com/bschaatsbergen/cidr/pulls) that you could help with.
