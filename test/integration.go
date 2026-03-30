package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chenjingxiong/weblise/server/protocol"
	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println("=== Testing Agent Connection ===")
	testAgentConnection("ws://localhost:10999/agent", "test-agent-123")

	fmt.Println()
	fmt.Println("=== Testing Client Connection ===")
	testClientConnection("ws://localhost:10999/client", "test-agent-123")
}

func testAgentConnection(url, key string) {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	defer ws.Close()

	// Send register message
	regMsg, _ := protocol.NewMessage(protocol.MessageTypeRegister, protocol.RegisterMessage{
		Key:    key,
		Name:   "Test Agent",
		OSType: "linux",
	})

	data, _ := json.Marshal(regMsg)
	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Fatalf("Write failed: %v", err)
	}

	fmt.Println("✓ Sent register message")

	// Send ping to verify connection
	pingMsg, _ := protocol.NewMessage(protocol.MessageTypeHeartbeat, protocol.HeartbeatMessage{
		Timestamp: time.Now().Unix(),
	})
	data, _ = json.Marshal(pingMsg)
	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Fatalf("Write ping failed: %v", err)
	}

	// Read response
	ws.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}

	var response protocol.Message
	if err := json.Unmarshal(msg, &response); err != nil {
		log.Printf("Raw response: %s", string(msg))
		return
	}

	if response.Type == protocol.MessageTypePong {
		fmt.Println("✓ Agent connected successfully (received pong)")
	} else {
		fmt.Printf("✓ Received: %s\n", response.Type)
	}

	// Keep connection alive briefly
	time.Sleep(500 * time.Millisecond)
	fmt.Println("✓ Agent connection test passed")
}

func testClientConnection(url, agentKey string) {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	defer ws.Close()

	// Send connect message
	connMsg, _ := protocol.NewMessage(protocol.MessageTypeConnect, protocol.ConnectMessage{
		AgentKey: agentKey,
	})

	data, _ := json.Marshal(connMsg)
	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Fatalf("Write failed: %v", err)
	}

	fmt.Println("✓ Sent connect message")

	// Read response
	ws.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}

	var response protocol.Message
	if err := json.Unmarshal(msg, &response); err != nil {
		log.Printf("Raw response: %s", string(msg))
		return
	}

	if response.Type == protocol.MessageTypeError {
		fmt.Println("✗ Agent not found (expected if agent disconnected)")
	} else if response.Type == protocol.MessageTypeConnect {
		fmt.Println("✓ Client connected to agent successfully!")
	} else {
		fmt.Printf("✓ Received: %s\n", response.Type)
	}
}
