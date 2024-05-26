package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type InterfaceConfig struct {
	InterfaceName string `json:"interfaceName"`
	MacAddress    string `json:"macAddress"`
	State         bool   `json:"state"`
	MTUSize       int    `json:"mtuSize"`
}

type NetworkDeviceConfig struct {
	Hostname   string                     `json:"hostname"`
	IPAddress  string                     `json:"ipAddress"`
	SubnetMask string                     `json:"subnetMask"`
	Gateway    string                     `json:"gateway"`
	DNSServers []string                   `json:"dnsServers"`
	Interfaces map[string]InterfaceConfig `json:"interfaces"`
}

func NewNetworkDeviceConfig(host string, ipAddr string, mask string, gw string, dns []string, ifaces map[string]InterfaceConfig) NetworkDeviceConfig {
	return NetworkDeviceConfig{Hostname: host, IPAddress: ipAddr, SubnetMask: mask, Gateway: gw, DNSServers: dns, Interfaces: ifaces}
}

func NewInterfaceConfig(ifName string, mac string, enabled bool, mtu int) InterfaceConfig {
	return InterfaceConfig{InterfaceName: ifName, MacAddress: mac, State: enabled, MTUSize: mtu}
}

func (device *NetworkDeviceConfig) AddInterface(ifName string, ifaceConfig InterfaceConfig) {
	device.Interfaces[ifName] = ifaceConfig
}

func (device *NetworkDeviceConfig) RemoveInterface(ifName string) {
	delete(device.Interfaces, ifName)
}

func (device *NetworkDeviceConfig) UpdateHostname(host string) {
	device.Hostname = host
}

func (device *NetworkDeviceConfig) UpdateIPAddress(ipAddr string) error {
	if net.ParseIP(ipAddr) != nil {
		device.IPAddress = ipAddr
		return nil
	} else {
		return fmt.Errorf("%s is not a valid IP address", ipAddr)
	}
}

func (device *NetworkDeviceConfig) UpdateSubnetMask(mask string) error {
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid subnet mask")
	}

	var combined uint32
	for _, part := range parts {
		octet, err := strconv.Atoi(part)
		if err != nil || octet < 0 || octet > 255 {
			return fmt.Errorf("Invalid subnet mask: %s", mask)
		}
		combined = (combined << 8) | uint32(octet)
	}
	maskBinStr := fmt.Sprintf("%32b", combined)
	if !strings.Contains(maskBinStr, "01") {
		return nil
	}
	return fmt.Errorf("Invalid subnet mask: %s", mask)
}

func (device *NetworkDeviceConfig) UpdateGateway(gw string) error {
	if net.ParseIP(gw) != nil {
		device.Gateway = gw
		return nil
	} else {
		return fmt.Errorf("%s is not a valid IP address", gw)
	}
}

func (device *NetworkDeviceConfig) UpdateDNSServers(dns []string) error {
	for _, server := range dns {
		if net.ParseIP(server) != nil {
			device.DNSServers = dns
			return nil
		}
	}
	return fmt.Errorf("DNS servers don't contain valid addresses")
}

func (device *NetworkDeviceConfig) UpdateInterface(ifName string, ifaceConfig InterfaceConfig) error {
	if _, present := device.Interfaces[ifName]; present {
		device.Interfaces[ifName] = ifaceConfig
		return nil
	} else {
		return fmt.Errorf("%s does not exist", ifName)
	}
}

func (device NetworkDeviceConfig) ShowDeviceConfig() {
	fmt.Printf("Hostname: %s\n", device.Hostname)
	fmt.Printf("IP: %s\n", device.IPAddress)
	fmt.Printf("SubnetMask: %s\n", device.SubnetMask)
	fmt.Printf("Gateway: %s\n", device.Gateway)
	fmt.Printf("DNS: %v\n", device.DNSServers)
	fmt.Printf("Interfaces: \n")
	for name, iface := range device.Interfaces {
		fmt.Printf("\tName: %s\n", name)
		fmt.Printf("\tMAC: %s\n", iface.MacAddress)
		fmt.Printf("\tState: %v\n", iface.State)
		fmt.Printf("\tMTU: %d\n", iface.MTUSize)
	}
}

func (device NetworkDeviceConfig) ToJSON() ([]byte, error) {
	d, err := json.Marshal(device)
	if err != nil {
		return nil, fmt.Errorf("Error serializing device configuration: %s", err)
	}
	return d, nil
}

func FromJSON(data []byte) ([]NetworkDeviceConfig, error) {
	var devices []NetworkDeviceConfig
	if err := json.Unmarshal(data, &devices); err != nil {
		return nil, fmt.Errorf("Error deserializing data %s", err)
	}
	return devices, nil
}

func main() {
	configs := flag.String("configs", "", "Add device configurations")
	flag.Parse()

	if *configs == "" {
		fmt.Printf("Usage: device-config [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	devices, err := FromJSON([]byte(*configs))
	if err != nil {
		fmt.Printf("Error importing device configs %s", err)
	}
	for _, device := range devices {
		device.ShowDeviceConfig()
		fmt.Println("-------------------------------")
	}
}
