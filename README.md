# tailhelp

A simple Go package that provides helper functions for working with Tailscale's API.

## Installation

```bash
go get github.com/yourusername/tailhelp
```

## Usage

```go
package main

import (
    "log"
    "github.com/yourusername/tailhelp"
)

func main() {
    // Get your Tailscale IPv4 address
    ip, err := tailhelp.MyIP()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("My Tailscale IP: %v", ip)

    // Find online peers with a specific hostname prefix
    peers, err := tailhelp.PeersFromHostnamePrefix("myprefix", true)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Found peers: %v", peers)
}
```

## Functions

- `MyIP() (netip.Addr, error)`: Returns the first IPv4 address of your Tailscale node
- `PeersFromHostnamePrefix(hostname string, onlyOnline bool) ([]string, error)`: Returns a list of peer DNS names that match the given hostname prefix

## Requirements

- Go 1.19 or later
- Tailscale daemon running on the system

## License

MIT License (or choose your preferred license)
