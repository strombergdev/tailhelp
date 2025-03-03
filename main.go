// Package tailhelp provides helper functions for working with Tailscale's API,
// making it easier to get information about the local node and its peers.
package tailhelp

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
	"sync"
	"time"

	"tailscale.com/client/tailscale"
	"tailscale.com/ipn/ipnstate"
)

var (
	// Global client instance with mutex for thread safety
	client     *tailscale.LocalClient
	clientOnce sync.Once
	clientMu   sync.RWMutex
)

// getClient returns the singleton LocalClient instance
func getClient() *tailscale.LocalClient {
	clientOnce.Do(func() {
		clientMu.Lock()
		defer clientMu.Unlock()
		client = &tailscale.LocalClient{}
	})

	clientMu.RLock()
	c := client
	clientMu.RUnlock()
	return c
}

// getStatus returns the current Tailscale node status with proper context management
func getStatus() (*ipnstate.Status, error) {
	// Create a context with timeout to ensure resources are released
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure context is canceled when function returns

	status, err := getClient().Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Tailscale status: %w", err)
	}
	return status, nil
}

// MyIP returns the first IPv4 address of the local Tailscale node.
// Returns an error if no IPv4 address is found or if unable to get node status.
func MyIP() (netip.Addr, error) {
	status, err := getStatus()
	if err != nil {
		return netip.Addr{}, err
	}

	// Find first IPv4 address
	for _, ip := range status.Self.TailscaleIPs {
		if ip.Is4() {
			return ip, nil
		}
	}

	return netip.Addr{}, fmt.Errorf("no IPv4 address found")
}

// PeersFromHostnamePrefix returns a list of peer DNS names that match the given hostname prefix.
// If onlyOnline is true, only returns peers that are currently online.
// Returns an empty list if no matching peers are found.
func PeersFromHostnamePrefix(hostname string, onlyOnline bool) ([]string, error) {
	status, err := getStatus()
	if err != nil {
		return nil, err
	}

	var peers []string
	for _, peer := range status.Peer {
		dnsName := strings.TrimSuffix(peer.DNSName, ".")
		if strings.HasPrefix(dnsName, hostname) {
			if onlyOnline && !peer.Online {
				continue
			}
			peers = append(peers, dnsName)
		}
	}

	return peers, nil
}

// PeerIPv4sFromHostnamePrefix returns a list of Tailscale IPv4 addresses for peers that match the given hostname prefix.
// If onlyOnline is true, only returns IPs from peers that are currently online.
// Returns an empty list if no matching peers are found.
func PeerIPv4sFromHostnamePrefix(hostname string, onlyOnline bool) ([]netip.Addr, error) {
	status, err := getStatus()
	if err != nil {
		return nil, err
	}

	var ips []netip.Addr
	for _, peer := range status.Peer {
		dnsName := strings.TrimSuffix(peer.DNSName, ".")
		if strings.HasPrefix(dnsName, hostname) {
			if onlyOnline && !peer.Online {
				continue
			}
			for _, ip := range peer.TailscaleIPs {
				if ip.Is4() {
					ips = append(ips, ip)
				}
			}
		}
	}

	return ips, nil
}
