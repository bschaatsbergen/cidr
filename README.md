# cidr
![release](https://github.com/bschaatsbergn/cidr/workflows/release/badge.svg?branch=main)

A cross platform CLI to perform various operations on a CIDR range.

## Brew
To install cidr using brew, simply do the below.

```sh
brew tap bschaatsbergen/cidr
brew install cidr
```

## Binaries
You can download the [latest binary](https://github.com/bschaatsbergn/cidr/releases/latest) for Linux, MacOS, and Windows.


## Usage

Foobar

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
65536
```

This also works with IPv6 addresses, for example:

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
