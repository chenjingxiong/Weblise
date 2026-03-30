package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chenjingxiong/weblise/server/config"
	"github.com/chenjingxiong/weblise/server/db"
	"github.com/chenjingxiong/weblise/server/http"
	"github.com/chenjingxiong/weblise/server/ws"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Weblise Server v0.1.0")

	cfg := config.Load()
	log.Printf("Config: HTTP=%s, WS=%s, DB=%s", cfg.HTTPPort, cfg.WSPort, cfg.DBPath)

	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	database, err := db.New(db.Config{Path: cfg.DBPath})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Println("Database initialized")

	hub := ws.NewHub()
	go hub.Run()
	log.Println("Hub started")

	wsServer := ws.NewServer(cfg.WSAddr(), hub)
	go func() {
		if err := wsServer.Start(); err != nil {
			log.Fatalf("WebSocket server failed: %v", err)
		}
	}()

	httpServer := http.NewServer(cfg.HTTPAddr())
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Println("Weblise Server is running")
	log.Printf("  - HTTP:  http://localhost%s", cfg.HTTPAddr())
	log.Printf("  - WebSocket: ws://localhost%s", cfg.WSAddr())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	httpServer.Shutdown(ctx)
	log.Println("Server stopped")
}
