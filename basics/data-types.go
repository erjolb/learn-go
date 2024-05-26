package main

import "fmt"

type Packet struct {
	sourceIP        string
	destinationIP   string
	sourcePort      uint16
	destinationPort uint16
	payload         string
}

// createPacket initializes a new Packet instance
func createPacket(srcIP string, dstIP string, srcPort uint16, dstPort uint16, data string) Packet {
	return Packet{sourceIP: srcIP, destinationIP: dstIP, sourcePort: srcPort, destinationPort: dstPort, payload: data}
}

func (p Packet) display() {
	fmt.Printf("Source IP: %s\n", p.sourceIP)
	fmt.Printf("Destination IP: %s\n", p.destinationIP)
	fmt.Printf("Source Port: %d\n", p.sourcePort)
	fmt.Printf("Destination Port: %d\n", p.destinationPort)
	fmt.Printf("Payload: %s\n", p.payload)
}

func main() {
	packet := createPacket("192.168.0.1", "192.168.0.5", 80, 65535, "Hello")

	packet.display()

}
