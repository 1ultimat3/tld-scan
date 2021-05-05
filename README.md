# tld-scan
Top level domain scanner in Go.

The core idea is to perform a TCP sync scan against potential web applications (running on port 80) for all possible top level domains.

## Installation
This tool can be installed using:
```bash
go build
```

## Usage
```bash
This program determines all used domains by checking if there is any web application running on <fqdm>:80.

Usage:
  tld scan [basename] [flags]

Flags:
  -h, --help          help for scan
  -w, --workers int   amount of workers running in parallel (default 10)
```

## Example
```bash
tld scan ethermat

ethermat.com
ethermat.ph
ethermat.pl
ethermat.vg
ethermat.ws
ethermat.xn--node
```
