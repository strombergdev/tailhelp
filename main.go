// Package tailhelp provides helper functions for working with Tailscale's API,
// making it easier to get information about the local node and its peers.
package tailhelp

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"tailscale.com/client/tailscale"
	"tailscale.com/ipn/ipnstate"
)

// getStatus returns the current Tailscale node status.
// It creates a new LocalClient and fetches the status from the Tailscale daemon.
func getStatus() (*ipnstate.Status, error) {
	client := &tailscale.LocalClient{}
	status, err := client.Status(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get Tailscale status: %v", err)
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
// Returns an error if no matching peers are found or if unable to get node status.
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

	if len(peers) == 0 {
		return nil, fmt.Errorf("no peer found with hostname: %s", hostname)
	}

	return peers, nil
}
