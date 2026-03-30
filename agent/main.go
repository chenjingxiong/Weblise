package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chenjingxiong/weblise/agent/conn"
)

var (
	serverAddr = flag.String("server", "", "WebSocket server address (ws://host:port)")
	deviceKey  = flag.String("key", "", "Device authentication key")
	name       = flag.String("name", "", "Device name (optional)")
	fps        = flag.Int("fps", 15, "Screen capture FPS (default: 15)")
	quality    = flag.Int("quality", 80, "JPEG quality 1-100 (default: 80)")
)

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Weblise Agent v0.1.0")

	if *serverAddr == "" {
		log.Fatal("--server is required (e.g., ws://localhost:8443)")
	}
	if *deviceKey == "" {
		log.Fatal("--key is required")
	}

	log.Printf("Server: %s", *serverAddr)
	log.Printf("Key: %s", *deviceKey)
	log.Printf("FPS: %d, Quality: %d", *fps, *quality)

	deviceName := *name
	if deviceName == "" {
		hostname, _ := os.Hostname()
		deviceName = hostname
	}

	connection, err := conn.New(&conn.Config{
		ServerAddr:     *serverAddr,
		DeviceKey:      *deviceKey,
		Name:          deviceName,
		CaptureFPS:    *fps,
		ScreenQuality: *quality,
	})
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}

	if err := connection.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	connection.Start()

	log.Println("Agent is running, press Ctrl+C to exit")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	connection.Close()
	log.Println("Agent stopped")
}
