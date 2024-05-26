package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// parseIPRange parses an IP range string and returns a slice of individual IP addresses.
func parseIPRange(ipRange string) []string {
	// Split the range into start and end IPs
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		fmt.Println("Invalid IP range format")
		return nil
	}

	startIP := parts[0]
	endIP := parts[1]

	start := ipToInt(startIP)
	end := ipToInt(endIP)

	if start == 0 || end == 0 || start > end {
		fmt.Println("Invalid IP range values")
		return nil
	}

	var ips []string
	for ip := start; ip <= end; ip++ {
		ips = append(ips, intToIP(ip))
	}

	return ips
}

// ipToInt converts an IP address string to a numeric format.
func ipToInt(ip string) uint32 {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0
	}

	var result uint32
	for _, part := range parts {
		var byteValue uint32
		fmt.Sscanf(part, "%d", &byteValue)
		result = result<<8 + byteValue
	}

	return result
}

// intToIP converts a numeric IP address back to its string format.
func intToIP(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(ip>>24)&0xFF,
		(ip>>16)&0xFF,
		(ip>>8)&0xFF,
		ip&0xFF)
}

// pingIP pings an IP address and sends the result to the results channel.
func pingIP(ip string, results chan<- string) {
	// Construct the ping command
	cmd := exec.Command("ping", "-c", "1", ip)

	// Run the ping command
	err := cmd.Run()

	// Check if the ping was successful
	if err != nil {
		results <- fmt.Sprintf("%s is down", ip)
	} else {
		results <- fmt.Sprintf("%s is up", ip)
	}
}

func main() {
	// Define command-line flags
	ipRange := flag.String("ip-range", "", "IP range to scan")

	// Parse the flags
	flag.Parse()

	// Check if the ip-range flag was provided
	if *ipRange == "" {
		fmt.Println("Usage: ipscanner [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Print the provided IP range
	fmt.Println("Scanning IP range:", *ipRange)

	// Parse the IP range into individual IP addresses
	ips := parseIPRange(*ipRange)
	if ips == nil {
		return
	}

	// Create a channel to collect the results
	results := make(chan string)

	// Launch goroutines to ping each IP address concurrently
	for _, ip := range ips {
		go pingIP(ip, results)
	}

	// Collect and print the results
	for i := 0; i < len(ips); i++ {
		fmt.Println(<-results)
	}
}
