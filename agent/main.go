package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	serverAddr = flag.String("server", "", "WebSocket server address (ws://host:port)")
	deviceKey  = flag.String("key", "", "Device authentication key")
)

func main() {
	flag.Parse()

	if *serverAddr == "" {
		log.Fatal("--server is required")
	}
	if *deviceKey == "" {
		log.Fatal("--key is required")
	}

	fmt.Printf("Weblise Agent v0.1.0\n")
	fmt.Printf("Server: %s\n", *serverAddr)
	fmt.Printf("Key: %s\n", *deviceKey)

	// TODO: Connect to server
}
