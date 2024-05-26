package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

var (
	clientCount uint8
	mu          sync.Mutex
)

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	clientCount++
	count := clientCount
	mu.Unlock()

	fmt.Fprintf(w, "Number of clients connected: %d\n", count)

}

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
