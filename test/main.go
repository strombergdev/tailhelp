package main

import (
	"fmt"
	"log"

	"github.com/strombergdev/tailhelp"
)

func main() {
	// Test getting peer DNS names
	peers, err := tailhelp.PeersFromHostnamePrefix("metrics-grafana", true)
	if err != nil {
		log.Fatalf("Failed to get peers: %v", err)
	}

	fmt.Println("Found peer DNS names:")
	for _, peer := range peers {
		fmt.Printf("- %s\n", peer)
	}

	// Test getting peer IPs
	ips, err := tailhelp.PeerIPv4sFromHostnamePrefix("metrics-grafana", true)
	if err != nil {
		log.Fatalf("Failed to get peer IPs: %v", err)
	}

	fmt.Println("\nFound peer IPs:")
	for _, ip := range ips {
		fmt.Printf("- %s\n", ip)
	}

	// Test getting local IP
	myIP, err := tailhelp.MyIP()
	if err != nil {
		log.Fatalf("Failed to get my IP: %v", err)
	}
	fmt.Printf("\nMy Tailscale IP: %s\n", myIP)
}
