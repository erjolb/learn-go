package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func readServerAcks(conn bufio.Reader) {
	message, err := conn.ReadString('\n')
	if err != nil {
		fmt.Printf("Error getting response from server %v\n", err)
	} else {
		fmt.Printf("Server: %v\n", message)
	}
}

func main() {
	// define the host
	host := "localhost:8080"

	// start the connection to the server
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Printf("Connection to the server failed with %v\n", err)
	}
	defer conn.Close()

	// get input from user
	scanner := bufio.NewScanner(os.Stdin)
	reader := bufio.NewReader(conn)

	for {
		fmt.Print("Enter message: ")
		if scanner.Scan() {
			input := scanner.Text()

			if input == "exit" {
				fmt.Println("Exiting...")
				break
			} else {
				conn.Write([]byte(input + "\n"))
			}
			readServerAcks(*reader)
		} else {
			fmt.Println("Error reading input: ", scanner.Err())
			break
		}
	}
}
