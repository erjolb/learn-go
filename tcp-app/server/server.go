package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	// Define the port on which the server will listen
	port := ":8080"

	// Start a TCP listener on the specified port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	// Print a message indicating that the server is listening
	fmt.Printf("Server is listening on port %s\n", port)

	// Loop to continuously accept new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Print a message indicating a new client connection
		fmt.Printf("New client connected: %s\n", conn.RemoteAddr().String())

		// Start a new Goroutine to handle the client
		go handleClient(conn)
	}
}

// handleClient handles communication with a single client
func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Read message from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			return
		}

		// Print the received message
		fmt.Printf("Received message from %s: %s", conn.RemoteAddr().String(), message)

		// Send a response back to the client
		response := "Message received\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Failed to send response: %v", err)
			return
		}
	}
}
