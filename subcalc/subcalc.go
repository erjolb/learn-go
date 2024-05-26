package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

func main() {
	// define command line arguments
	subnet := flag.String("subnet", "", "enter subnet in CIDR notation. e.g. 192.168.0.1/24")
	flag.Parse()

	if *subnet == "" {
		fmt.Println("Invalid subnet format. Use CIDR notation e.g. 192.168.0.1/24")
		return
	}

	// split ip and subnet
	parts := strings.Split(*subnet, "/")
	if len(parts) != 2 {
		fmt.Println("Error: Invalid CIDR format. Make sure you include both IP and mask.")
		return
	}

	ipv4 := net.ParseIP(parts[0])
	if ipv4 == nil {
		fmt.Println("Invalid IP address format")
		return
	}

	fmt.Printf("IPv4: %v\n", ipv4)

}
