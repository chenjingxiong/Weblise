# Weblise MVP Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a minimal viable remote desktop system with Agent (Go single binary), Server (Go + Docker), and Web Client (HTML/JS) that can display remote screen and control mouse/keyboard input.

**Architecture:**
- **Agent**: Go single binary that captures screen and executes input commands, connects to Server via WebSocket
- **Server**: Go HTTP/WebSocket server with SQLite for device management, routes messages between Client and Agent
- **Web Client**: Single HTML page with JavaScript for WebSocket communication and screen display
- **Deployment**: Docker container with two ports (8080 for HTTP, 8443 for WebSocket), works behind reverse proxy

**Tech Stack:** Go 1.21+, SQLite, vanilla JavaScript, WebSocket

---

## File Structure

```
weblise/
├── agent/                      # Agent (Go single binary)
│   ├── main.go                 # Entry point, command-line flags
│   ├── screen/
│   │   ├── capture.go          # Screen capture interface
│   │   ├── windows.go          # Windows implementation
│   │   └── linux.go            # Linux implementation
│   ├── input/
│   │   ├── mouse.go            # Mouse control
│   │   └── keyboard.go         # Keyboard control
│   └── conn/
│       └── ws.go               # WebSocket connection to server
│
├── server/                     # Server (Go)
│   ├── main.go                 # Entry point
│   ├── cmd/
│   │   └── root.go             # CLI command setup
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── http/
│   │   ├── server.go           # HTTP server (:8080)
│   │   ├── handler.go          # HTTP handlers
│   │   └── static.go           # Embedded static files
│   ├── ws/
│   │   ├── server.go           # WebSocket server (:8443)
│   │   ├── hub.go              # Connection hub and routing
│   │   ├── agent.go            # Agent connection handler
│   │   └── client.go           # Client connection handler
│   ├── db/
│   │   ├── db.go               # SQLite connection
│   │   └── schema.sql          # Database schema
│   └── protocol/
│       └── messages.go         # Message type definitions
│
├── web/                        # Web client (embedded in server)
│   ├── index.html              # Single page app
│   ├── app.js                  # Main app logic
│   ├── screen.js               # Screen rendering
│   ├── input.js                # Input capture and send
│   └── style.css               # Styles
│
├── docker/
│   ├── Dockerfile              # Server Docker image
│   └── docker-compose.yml      # Docker compose config
│
├── scripts/
│   ├── build-agent.sh          # Build agent binary
│   └── build-server.sh         # Build server binary
│
├── docs/
│   └── superpowers/
│       ├── specs/
│       │   └── 2026-03-30-weblise-mvp-design.md
│       └── plans/
│           └── 2026-03-30-weblise-mvp-implementation.md
│
├── go.mod                      # Go module definition
└── README.md                   # Project documentation
```

---

## Chunk 1: Project Foundation

### Task 1: Initialize Go Module and Project Structure

**Files:**
- Create: `go.mod`
- Create: `README.md`
- Create: `agent/main.go`
- Create: `server/main.go`
- Create: `web/index.html` (placeholder)

- [ ] **Step 1: Initialize Go module**

```bash
cd /root/projects/Weblise
go mod init github.com/chenjingxiong/weblise
```

Expected: `go: creating new go.mod: module github.com/chenjingxiong/weblise`

- [ ] **Step 2: Create go.mod with dependencies**

```go
// go.mod
module github.com/chenjingxiong/weblise

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gorilla/websocket v1.5.1
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/kbinani/screenshot v0.0.0-20230601120736-8b421e4c6696
	github.com/go-vgo/robotgo v0.0.0-20230830182337-9b5c1e6ab2f6
	github.com/spf13/cobra v1.8.0
	github.com/shirou/gopsutil/v3 v3.23.12
)
```

- [ ] **Step 3: Download dependencies**

```bash
go mod tidy
```

Expected: Downloads all dependencies

- [ ] **Step 4: Create README.md**

```markdown
# Weblise

Web-based remote desktop system - MVP

## Quick Start

### Server

```bash
cd server
go run main.go serve
```

Server runs on:
- HTTP: http://localhost:8080
- WebSocket: ws://localhost:8443

### Agent

```bash
cd agent
go run main.go --server=ws://localhost:8443 --key=test-device-key
```

## Architecture

- **Agent**: Go binary for screen capture and input control
- **Server**: Go server for routing messages between clients and agents
- **Web Client**: Browser-based remote control interface

## Development Status

- [x] Design specification
- [ ] Foundation
- [ ] Agent development
- [ ] Server development
- [ ] Web client development
- [ ] Docker deployment
- [ ] End-to-end testing
```

- [ ] **Step 5: Create agent main.go stub**

```go
// agent/main.go
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
```

- [ ] **Step 6: Create server main.go stub**

```go
// server/main.go
package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Printf("Weblise Server v0.1.0\n")

	// TODO: Start HTTP server
	// TODO: Start WebSocket server
}
```

- [ ] **Step 7: Create web index.html placeholder**

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weblise - Remote Desktop</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <div class="container">
        <h1>Weblise Remote Desktop</h1>
        <div id="connection-panel">
            <input type="text" id="agent-key" placeholder="Enter Agent Key">
            <button id="connect-btn">Connect</button>
        </div>
        <div id="status" class="status">Disconnected</div>
        <div id="screen-container" style="display: none;">
            <canvas id="screen-canvas"></canvas>
        </div>
    </div>
    <script src="/app.js"></script>
</body>
</html>
```

- [ ] **Step 8: Create placeholder files**

```bash
touch web/app.js web/screen.js web/input.js web/style.css
```

- [ ] **Step 9: Test agent compiles**

```bash
cd /root/projects/Weblise/agent
go build -o agent main.go
```

Expected: Creates `agent` binary without errors

- [ ] **Step 10: Test server compiles**

```bash
cd /root/projects/Weblise/server
go build -o server main.go
```

Expected: Creates `server` binary without errors

- [ ] **Step 11: Commit foundation**

```bash
cd /root/projects/Weblise
git add go.mod go.sum README.md agent/main.go server/main.go web/
git commit -m "feat: initialize project structure and Go module

- Add go.mod with dependencies
- Create README with project overview
- Add agent and server main.go stubs
- Add web client placeholder files

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 2: Protocol and Data Types

### Task 2: Define Protocol Message Types

**Files:**
- Create: `server/protocol/messages.go`
- Create: `server/protocol/types.go`

- [ ] **Step 1: Create protocol/messages.go with all message types**

```go
// server/protocol/messages.go
package protocol

import "encoding/json"

// MessageType represents the type of message
type MessageType string

const (
	// Agent messages
	MessageTypeRegister     MessageType = "register"
	MessageTypeHeartbeat    MessageType = "ping"
	MessageTypeError        MessageType = "error"
	MessageTypeFrame        MessageType = "frame"

	// Client messages
	MessageTypeConnect      MessageType = "connect"
	MessageTypeDisconnect   MessageType = "disconnect"
	MessageTypeInput        MessageType = "input"

	// Bidirectional
	MessageTypePong         MessageType = "pong"
)

// Message is the base message structure
type Message struct {
	Type MessageType          `json:"type"`
	Data json.RawMessage      `json:"data,omitempty"`
}

// RegisterMessage is sent by agent to register with server
type RegisterMessage struct {
	Key   string `json:"key"`
	Name  string `json:"name,omitempty"`
	OSType string `json:"os_type,omitempty"`
}

// ConnectMessage is sent by client to connect to an agent
type ConnectMessage struct {
	AgentKey string `json:"agent_key"`
}

// FrameMessage contains screen frame data
type FrameMessage struct {
	Data      []byte `json:"data"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Timestamp int64  `json:"timestamp"`
}

// InputMessage contains input event data
type InputMessage struct {
	Action string      `json:"action"` // mousemove, mousedown, mouseup, keydown, keyup
	Data   InputData   `json:"data"`
}

// InputData represents input event data
type InputData struct {
	X       int    `json:"x,omitempty"`
	Y       int    `json:"y,omitempty"`
	Button  int    `json:"button,omitempty"`
	Key     string `json:"key,omitempty"`
	KeyCode int    `json:"keyCode,omitempty"`
}

// HeartbeatMessage for keep-alive
type HeartbeatMessage struct {
	Timestamp int64 `json:"timestamp"`
}

// ErrorMessage for error reporting
type ErrorMessage struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// NewMessage creates a new message with data
func NewMessage(t MessageType, data interface{}) (*Message, error) {
	var rawData json.RawMessage
	if data != nil {
		bytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		rawData = bytes
	}
	return &Message{Type: t, Data: rawData}, nil
}

// ParseData parses message data into target struct
func (m *Message) ParseData(target interface{}) error {
	if m.Data == nil {
		return nil
	}
	return json.Unmarshal(m.Data, target)
}
```

- [ ] **Step 2: Create protocol/types.go for connection types**

```go
// server/protocol/types.go
package protocol

import (
	"time"
	"github.com/google/uuid"
)

// ConnectionType represents the type of connection
type ConnectionType string

const (
	ConnectionTypeAgent  ConnectionType = "agent"
	ConnectionTypeClient ConnectionType = "client"
)

// Connection represents a WebSocket connection
type Connection struct {
	ID           string
	Type         ConnectionType
	DeviceKey    string
	ConnectedAt  time.Time
	LastSeen     time.Time
	SendChan     chan *Message
	CloseChan    chan struct{}
}

// NewConnection creates a new connection
func NewConnection(connType ConnectionType, deviceKey string) *Connection {
	return &Connection{
		ID:          uuid.New().String(),
		Type:        connType,
		DeviceKey:   deviceKey,
		ConnectedAt: time.Now(),
		LastSeen:    time.Now(),
		SendChan:    make(chan *Message, 256),
		CloseChan:   make(chan struct{}),
	}
}

// Device represents a registered device
type Device struct {
	ID        string
	Key       string
	Name      string
	OSType    string
	LastSeen  time.Time
	CreatedAt time.Time
}

// Session represents an active remote desktop session
type Session struct {
	ID           string
	DeviceID     string
	ClientConnID string
	AgentConnID  string
	StartedAt    time.Time
}
```

- [ ] **Step 3: Add uuid dependency**

```bash
cd /root/projects/Weblise
go get github.com/google/uuid
```

Expected: Adds uuid dependency

- [ ] **Step 4: Test protocol package compiles**

```bash
cd /root/projects/Weblise
go build ./server/protocol/...
```

Expected: Compiles without errors

- [ ] **Step 5: Create protocol test**

```go
// server/protocol/messages_test.go
package protocol

import (
	"encoding/json"
	"testing"
)

func TestNewMessage(t *testing.T) {
	data := RegisterMessage{Key: "test-key", Name: "Test Device"}
	msg, err := NewMessage(MessageTypeRegister, data)
	if err != nil {
		t.Fatalf("NewMessage failed: %v", err)
	}

	if msg.Type != MessageTypeRegister {
		t.Errorf("Expected type %s, got %s", MessageTypeRegister, msg.Type)
	}

	var parsed RegisterMessage
	if err := msg.ParseData(&parsed); err != nil {
		t.Fatalf("ParseData failed: %v", err)
	}

	if parsed.Key != "test-key" {
		t.Errorf("Expected key 'test-key', got '%s'", parsed.Key)
	}
}

func TestInputMessage(t *testing.T) {
	msg := InputMessage{
		Action: "mousemove",
		Data: InputData{
			X: 100,
			Y: 200,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded InputMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Action != "mousemove" {
		t.Errorf("Expected action 'mousemove', got '%s'", decoded.Action)
	}

	if decoded.Data.X != 100 || decoded.Data.Y != 200 {
		t.Errorf("Expected coordinates (100, 200), got (%d, %d)", decoded.Data.X, decoded.Data.Y)
	}
}
```

- [ ] **Step 6: Run tests**

```bash
cd /root/projects/Weblise
go test ./server/protocol/... -v
```

Expected: All tests pass

- [ ] **Step 7: Commit protocol implementation**

```bash
cd /root/projects/Weblise
git add server/protocol/ go.mod go.sum
git commit -m "feat: implement protocol message types

- Define message types for agent/client communication
- Add connection and session types
- Add unit tests for message parsing
- Add uuid dependency for connection IDs

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 3: Database Layer

### Task 3: Implement SQLite Database Layer

**Files:**
- Create: `server/db/schema.sql`
- Create: `server/db/db.go`

- [ ] **Step 1: Create database schema**

```sql
-- server/db/schema.sql

-- Devices table
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    name TEXT,
    os_type TEXT,
    last_seen INTEGER,
    created_at INTEGER
);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    device_id TEXT,
    client_conn_id TEXT,
    agent_conn_id TEXT,
    started_at INTEGER,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);

-- Index for faster device lookups by key
CREATE INDEX IF NOT EXISTS idx_devices_key ON devices(key);
```

- [ ] **Step 2: Create database layer**

```go
// server/db/db.go
package db

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var SchemaFS embed.FS

type Database struct {
	db *sql.DB
}

// Config holds database configuration
type Config struct {
	Path string // Path to SQLite database file
}

// New creates a new database connection
func New(cfg Config) (*Database, error) {
	db, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	// Run schema
	if err := runSchema(db); err != nil {
		return nil, fmt.Errorf("run schema: %w", err)
	}

	return &Database{db: db}, nil
}

// runSchema executes the database schema
func runSchema(db *sql.DB) error {
	schema, err := SchemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	_, err = db.Exec(string(schema))
	return err
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// Device operations

// CreateDevice creates a new device
func (d *Database) CreateDevice(id, key, name, osType string) error {
	now := time.Now().Unix()
	query := `
		INSERT INTO devices (id, key, name, os_type, last_seen, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(query, id, key, name, osType, now, now)
	return err
}

// GetDeviceByKey retrieves a device by its key
func (d *Database) GetDeviceByKey(key string) (*Device, error) {
	query := `
		SELECT id, key, name, os_type, last_seen, created_at
		FROM devices WHERE key = ?
	`
	row := d.db.QueryRow(query, key)

	var d Device
	var lastSeen, createdAt int64
	err := row.Scan(&d.ID, &d.Key, &d.Name, &d.OSType, &lastSeen, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	d.LastSeen = time.Unix(lastSeen, 0)
	d.CreatedAt = time.Unix(createdAt, 0)
	return &d, nil
}

// UpdateDeviceLastSeen updates the last seen timestamp for a device
func (d *Database) UpdateDeviceLastSeen(id string) error {
	query := `UPDATE devices SET last_seen = ? WHERE id = ?`
	_, err := d.db.Exec(query, time.Now().Unix(), id)
	return err
}

// Session operations

// CreateSession creates a new session
func (d *Database) CreateSession(id, deviceID, clientConnID, agentConnID string) error {
	query := `
		INSERT INTO sessions (id, device_id, client_conn_id, agent_conn_id, started_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(query, id, deviceID, clientConnID, agentConnID, time.Now().Unix())
	return err
}

// DeleteSession deletes a session by ID
func (d *Database) DeleteSession(id string) error {
	_, err := d.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

// Device represents a device record
type Device struct {
	ID        string
	Key       string
	Name      string
	OSType    string
	LastSeen  time.Time
	CreatedAt time.Time
}
```

- [ ] **Step 3: Create database tests**

```go
// server/db/db_test.go
package db

import (
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Use in-memory database for testing
	cfg := Config{Path: ":memory:"}
	db, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer db.Close()

	// Check if tables were created
	var tableName string
	err = db.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='devices'").Scan(&tableName)
	if err != nil {
		t.Errorf("Devices table not created: %v", err)
	}
}

func TestCreateDevice(t *testing.T) {
	cfg := Config{Path":memory:"}
	db, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer db.Close()

	err = db.CreateDevice("device-1", "test-key", "Test Device", "linux")
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}

	device, err := db.GetDeviceByKey("test-key")
	if err != nil {
		t.Fatalf("GetDeviceByKey failed: %v", err)
	}

	if device == nil {
		t.Fatal("Device not found")
	}

	if device.Key != "test-key" {
		t.Errorf("Expected key 'test-key', got '%s'", device.Key)
	}

	if device.Name != "Test Device" {
		t.Errorf("Expected name 'Test Device', got '%s'", device.Name)
	}
}

func TestGetDeviceByKeyNotFound(t *testing.T) {
	cfg := Config{Path: ":memory:"}
	db, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer db.Close()

	device, err := db.GetDeviceByKey("non-existent")
	if err != nil {
		t.Fatalf("GetDeviceByKey failed: %v", err)
	}

	if device != nil {
		t.Error("Expected nil for non-existent device")
	}
}

func TestUpdateDeviceLastSeen(t *testing.T) {
	cfg := Config{Path: ":memory:"}
	db, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer db.Close()

	err = db.CreateDevice("device-1", "test-key", "Test Device", "linux")
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}

	// Wait a bit to ensure timestamp changes
	time.Sleep(10 * time.Millisecond)

	err = db.UpdateDeviceLastSeen("device-1")
	if err != nil {
		t.Fatalf("UpdateDeviceLastSeen failed: %v", err)
	}

	device, err := db.GetDeviceByKey("test-key")
	if err != nil {
		t.Fatalf("GetDeviceByKey failed: %v", err)
	}

	// Check that last seen was updated (should be very recent)
	if time.Since(device.LastSeen) > 5*time.Second {
		t.Error("LastSeen was not updated properly")
	}
}
```

- [ ] **Step 4: Fix test syntax error**

```go
// Fix the TestNew function - missing = in Config
cfg := Config{Path: ":memory:"}
```

- [ ] **Step 5: Run database tests**

```bash
cd /root/projects/Weblise
go test ./server/db/... -v
```

Expected: All tests pass

- [ ] **Step 6: Commit database layer**

```bash
cd /root/projects/Weblise
git add server/db/
git commit -m "feat: implement SQLite database layer

- Add schema for devices and sessions
- Implement database operations (CRUD)
- Add unit tests for database operations
- Use embed for schema file

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 4: WebSocket Hub and Routing

### Task 4: Implement WebSocket Connection Hub

**Files:**
- Create: `server/ws/hub.go`

- [ ] **Step 1: Create WebSocket hub**

```go
// server/ws/hub.go
package ws

import (
	"log"
	"sync"

	"github.com/chenjingxiong/weblise/server/protocol"
)

// Hub manages active connections and routes messages between agents and clients
type Hub struct {
	mu            sync.RWMutex
	agents        map[string]*Connection  // key: connection ID
	clients       map[string]*Connection  // key: connection ID
	agentsByKey   map[string]*Connection  // key: device key
	sessions      map[string]*Session     // key: session ID
	register      chan *Connection
	unregister    chan *Connection
	messageAgent  chan *AgentMessage
	messageClient chan *ClientMessage
}

// AgentMessage represents a message from an agent
type AgentMessage struct {
	ConnID string
	Msg    *protocol.Message
}

// ClientMessage represents a message from a client
type ClientMessage struct {
	ConnID string
	Msg    *protocol.Message
}

// Session represents an active remote desktop session
type Session struct {
	ID         string
	AgentConn  *Connection
	ClientConn *Connection
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		agents:        make(map[string]*Connection),
		clients:       make(map[string]*Connection),
		agentsByKey:   make(map[string]*Connection),
		sessions:      make(map[string]*Session),
		register:      make(chan *Connection),
		unregister:    make(chan *Connection),
		messageAgent:  make(chan *AgentMessage, 256),
		messageClient: make(chan *ClientMessage, 256),
	}
}

// Run starts the hub's event loop
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.handleRegister(conn)

		case conn := <-h.unregister:
			h.handleUnregister(conn)

		case msg := <-h.messageAgent:
			h.handleAgentMessage(msg)

		case msg := <-h.messageClient:
			h.handleClientMessage(msg)
		}
	}
}

// handleRegister registers a new connection
func (h *Hub) handleRegister(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("[Hub] Registering %s connection: %s (key: %s)", conn.Type, conn.ID, conn.DeviceKey)

	switch conn.Type {
	case protocol.ConnectionTypeAgent:
		h.agents[conn.ID] = conn
		h.agentsByKey[conn.DeviceKey] = conn
		log.Printf("[Hub] Agent registered: %s", conn.DeviceKey)

	case protocol.ConnectionTypeClient:
		h.clients[conn.ID] = conn
		log.Printf("[Hub] Client registered: %s", conn.ID)
	}
}

// handleUnregister unregisters a connection
func (h *Hub) handleUnregister(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("[Hub] Unregistering %s connection: %s", conn.Type, conn.ID)

	switch conn.Type {
	case protocol.ConnectionTypeAgent:
		delete(h.agents, conn.ID)
		delete(h.agentsByKey, conn.DeviceKey)

		// Clean up any sessions
		for sessionID, session := range h.sessions {
			if session.AgentConn.ID == conn.ID {
				delete(h.sessions, sessionID)
				// Notify client
				if session.ClientConn != nil {
					session.ClientConn.SendSafe(protocol.NewErrorMessage("Agent disconnected"))
				}
			}
		}
		log.Printf("[Hub] Agent unregistered: %s", conn.DeviceKey)

	case protocol.ConnectionTypeClient:
		delete(h.clients, conn.ID)

		// Clean up session
		for sessionID, session := range h.sessions {
			if session.ClientConn.ID == conn.ID {
				delete(h.sessions, sessionID)
			}
		}
		log.Printf("[Hub] Client unregistered: %s", conn.ID)
	}
}

// handleAgentMessage handles a message from an agent
func (h *Hub) handleAgentMessage(msg *AgentMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Find the session for this agent
	for _, session := range h.sessions {
		if session.AgentConn.ID == msg.ConnID {
			// Forward to client
			session.ClientConn.SendSafe(msg.Msg)
			return
		}
	}

	log.Printf("[Hub] No session found for agent message: %s", msg.ConnID)
}

// handleClientMessage handles a message from a client
func (h *Hub) handleClientMessage(msg *ClientMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Find the session for this client
	for _, session := range h.sessions {
		if session.ClientConn.ID == msg.ConnID {
			// Forward to agent
			session.AgentConn.SendSafe(msg.Msg)
			return
		}
	}

	log.Printf("[Hub] No session found for client message: %s", msg.ConnID)
}

// RegisterAgentConnection registers an agent connection
func (h *Hub) RegisterAgentConnection(conn *Connection) {
	h.register <- conn
}

// RegisterClientConnection registers a client connection
func (h *Hub) RegisterClientConnection(conn *Connection) {
	h.register <- conn
}

// Unregister unregisters a connection
func (h *Hub) Unregister(conn *Connection) {
	h.unregister <- conn
}

// GetAgentByKey retrieves an agent connection by device key
func (h *Hub) GetAgentByKey(key string) *Connection {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.agentsByKey[key]
}

// CreateSession creates a new session between client and agent
func (h *Hub) CreateSession(sessionID string, clientConn, agentConn *Connection) *Session {
	h.mu.Lock()
	defer h.mu.Unlock()

	session := &Session{
		ID:         sessionID,
		AgentConn:  agentConn,
		ClientConn: clientConn,
	}
	h.sessions[sessionID] = session

	log.Printf("[Hub] Session created: %s (client=%s, agent=%s)",
		sessionID, clientConn.ID, agentConn.ID)

	return session
}

// DeleteSession deletes a session
func (h *Hub) DeleteSession(sessionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.sessions, sessionID)
	log.Printf("[Hub] Session deleted: %s", sessionID)
}

// NewErrorMessage creates an error message
func NewErrorMessage(message string) *protocol.Message {
	return &protocol.Message{
		Type: protocol.MessageTypeError,
		Data: nil, // Simplified for MVP
	}
}
```

- [ ] **Step 2: Create Connection wrapper with send utilities**

```go
// server/ws/connection.go
package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/chenjingxiong/weblise/server/protocol"
)

// Connection wraps a WebSocket connection with protocol message handling
type Connection struct {
	*protocol.Connection
	mu sync.Mutex
}

// NewConnection creates a new connection wrapper
func NewConnection(base *protocol.Connection) *Connection {
	return &Connection{Connection: base}
}

// Send sends a message to the connection
func (c *Connection) Send(msg *protocol.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.SendChan <- msg:
		return nil
	case <-c.CloseChan:
		return nil
	}
}

// SendSafe sends a message without blocking (drops if full)
func (c *Connection) SendSafe(msg *protocol.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.SendChan <- msg:
	default:
		log.Printf("[Connection] Send channel full, dropping message: %s", c.ID)
	}
}

// SendJSON sends a JSON message
func (c *Connection) SendJSON(t protocol.MessageType, data interface{}) error {
	msg, err := protocol.NewMessage(t, data)
	if err != nil {
		return err
	}
	return c.Send(msg)
}

// SendBytes sends raw bytes (for screen frames)
func (c *Connection) SendBytes(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.SendChan <- &protocol.Message{Type: protocol.MessageTypeFrame}:
		// For binary data, we'll handle this differently in the actual WebSocket handler
		return nil
	case <-c.CloseChan:
		return nil
	}
}

// MarshalMessage converts a message to JSON
func MarshalMessage(msg *protocol.Message) ([]byte, error) {
	return json.Marshal(msg)
}

// UnmarshalMessage parses JSON into a message
func UnmarshalMessage(data []byte) (*protocol.Message, error) {
	var msg protocol.Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
```

- [ ] **Step 3: Create hub tests**

```go
// server/ws/hub_test.go
package ws

import (
	"testing"
	"time"

	"github.com/chenjingxiong/weblise/server/protocol"
	"github.com/google/uuid"
)

func TestHub_RegisterAgentConnection(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop() // We'll add Stop method

	conn := NewConnection(protocol.NewConnection(protocol.ConnectionTypeAgent, "test-agent-key"))
	hub.RegisterAgentConnection(conn.Connection)

	time.Sleep(50 * time.Millisecond) // Let the hub process

	agent := hub.GetAgentByKey("test-agent-key")
	if agent == nil {
		t.Fatal("Agent not found after registration")
	}

	if agent.DeviceKey != "test-agent-key" {
		t.Errorf("Expected device key 'test-agent-key', got '%s'", agent.DeviceKey)
	}
}

func TestHub_CreateSession(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	agentConn := NewConnection(protocol.NewConnection(protocol.ConnectionTypeAgent, "agent-key"))
	clientConn := NewConnection(protocol.NewConnection(protocol.ConnectionTypeClient, ""))

	hub.RegisterAgentConnection(agentConn.Connection)
	hub.RegisterClientConnection(clientConn.Connection)

	time.Sleep(50 * time.Millisecond)

	sessionID := uuid.New().String()
	session := hub.CreateSession(sessionID, clientConn.Connection, agentConn.Connection)

	if session == nil {
		t.Fatal("Session not created")
	}

	if session.ID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, session.ID)
	}

	if session.AgentConn.ID != agentConn.ID {
		t.Error("Agent connection mismatch")
	}

	if session.ClientConn.ID != clientConn.ID {
		t.Error("Client connection mismatch")
	}
}
```

- [ ] **Step 4: Add Stop method to Hub**

```go
// Add to hub.go

// Stop gracefully stops the hub
func (h *Hub) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Close all connections
	for _, conn := range h.agents {
		close(conn.SendChan)
		close(conn.CloseChan)
	}
	for _, conn := range h.clients {
		close(conn.SendChan)
		close(conn.CloseChan)
	}
}
```

- [ ] **Step 5: Run tests**

```bash
cd /root/projects/Weblise
go test ./server/ws/... -v
```

Expected: Tests pass (may need adjustments)

- [ ] **Step 6: Commit WebSocket hub**

```bash
cd /root/projects/Weblise
git add server/ws/
git commit -m "feat: implement WebSocket hub for connection management

- Add hub for managing agent and client connections
- Implement message routing between agents and clients
- Add session management
- Add connection utilities for safe message sending
- Add unit tests

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 5: WebSocket Server Implementation

### Task 5: Implement WebSocket Server

**Files:**
- Create: `server/ws/server.go`
- Create: `server/ws/agent.go`
- Create: `server/ws/client.go`

- [ ] **Step 1: Create WebSocket server**

```go
// server/ws/server.go
package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 32, // 32KB for screen frames
	WriteBufferSize: 1024 * 32,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for MVP
	},
}

// Server manages the WebSocket server
type Server struct {
	addr    string
	hub     *Hub
	agentMux sync.Mutex
	clientMux sync.Mutex
}

// NewServer creates a new WebSocket server
func NewServer(addr string, hub *Hub) *Server {
	return &Server{
		addr: addr,
		hub:  hub,
	}
}

// Start starts the WebSocket server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Agent endpoint
	mux.HandleFunc("/agent", s.handleAgentConnection)

	// Client endpoint
	mux.HandleFunc("/client", s.handleClientConnection)

	log.Printf("[WS Server] Starting on %s", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

// handleAgentConnection handles incoming agent connections
func (s *Server) handleAgentConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS Server] Agent upgrade failed: %v", err)
		return
	}

	log.Printf("[WS Server] New agent connection from %s", r.RemoteAddr)
	go s.handleAgent(conn)
}

// handleClientConnection handles incoming client connections
func (s *Server) handleClientConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS Server] Client upgrade failed: %v", err)
		return
	}

	log.Printf("[WS Server] New client connection from %s", r.RemoteAddr)
	go s.handleClient(conn)
}
```

- [ ] **Step 2: Create agent connection handler**

```go
// server/ws/agent.go
package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/chenjingxiong/weblise/server/db"
	"github.com/chenjingxiong/weblise/server/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// handleAgent handles an agent WebSocket connection
func (s *Server) handleAgent(wsConn *websocket.Conn) {
	defer wsConn.Close()

	// Wait for register message
	var msg protocol.Message
	if err := wsConn.ReadJSON(&msg); err != nil {
		log.Printf("[Agent] Read register message failed: %v", err)
		return
	}

	if msg.Type != protocol.MessageTypeRegister {
		log.Printf("[Agent] First message must be register, got: %s", msg.Type)
		return
	}

	var regMsg protocol.RegisterMessage
	if err := msg.ParseData(&regMsg); err != nil {
		log.Printf("[Agent] Parse register message failed: %v", err)
		return
	}

	log.Printf("[Agent] Register: key=%s, name=%s, os=%s", regMsg.Key, regMsg.Name, regMsg.OSType)

	// Create connection
	baseConn := protocol.NewConnection(protocol.ConnectionTypeAgent, regMsg.Key)
	conn := NewConnection(baseConn)

	// Register with hub
	s.hub.RegisterAgentConnection(conn.Connection)

	// Start send goroutine
	done := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-conn.SendChan:
				data, err := json.Marshal(msg)
				if err != nil {
					log.Printf("[Agent] Marshal message failed: %v", err)
					continue
				}
				if err := wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("[Agent] Send message failed: %v", err)
					close(done)
					return
				}
			case <-conn.CloseChan:
				return
			case <-done:
				return
			}
		}
	}()

	// Read loop
	for {
		messageType, data, err := wsConn.ReadMessage()
		if err != nil {
			log.Printf("[Agent] Read message failed: %v", err)
			break
		}

		conn.LastSeen = time.Now()

		if messageType == websocket.TextMessage {
			var msg protocol.Message
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("[Agent] Unmarshal message failed: %v", err)
				continue
			}

			// Handle message types
			switch msg.Type {
			case protocol.MessageTypeHeartbeat:
				// Respond with pong
				conn.SendSafe(&protocol.Message{Type: protocol.MessageTypePong})

			default:
				// Forward to hub
				s.hub.messageAgent <- &AgentMessage{
					ConnID: conn.ID,
					Msg:    &msg,
				}
			}
		} else if messageType == websocket.BinaryMessage {
			// Binary data (screen frame)
			// For MVP, we'll send as base64 in JSON
			// Later: optimize with pure binary
		}
	}

	// Cleanup
	close(done)
	s.hub.Unregister(conn.Connection)
	log.Printf("[Agent] Connection closed: %s", conn.ID)
}
```

- [ ] **Step 3: Create client connection handler**

```go
// server/ws/client.go
package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/chenjingxiong/weblise/server/db"
	"github.com/chenjingxiong/weblise/server/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// handleClient handles a client WebSocket connection
func (s *Server) handleClient(wsConn *websocket.Conn) {
	defer wsConn.Close()

	// Create connection (no key initially)
	baseConn := protocol.NewConnection(protocol.ConnectionTypeClient, "")
	conn := NewConnection(baseConn)

	// Register with hub
	s.hub.RegisterClientConnection(conn.Connection)

	// Start send goroutine
	done := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-conn.SendChan:
				data, err := json.Marshal(msg)
				if err != nil {
					log.Printf("[Client] Marshal message failed: %v", err)
					continue
				}
				if err := wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("[Client] Send message failed: %v", err)
					close(done)
					return
				}
			case <-conn.CloseChan:
				return
			case <-done:
				return
			}
		}
	}()

	// Read loop
	for {
		messageType, data, err := wsConn.ReadMessage()
		if err != nil {
			log.Printf("[Client] Read message failed: %v", err)
			break
		}

		conn.LastSeen = time.Now()

		if messageType == websocket.TextMessage {
			var msg protocol.Message
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("[Client] Unmarshal message failed: %v", err)
				continue
			}

			// Handle message types
			switch msg.Type {
			case protocol.MessageTypeConnect:
				var connMsg protocol.ConnectMessage
				if err := msg.ParseData(&connMsg); err != nil {
					log.Printf("[Client] Parse connect message failed: %v", err)
					continue
				}

				// Find agent
				agentConn := s.hub.GetAgentByKey(connMsg.AgentKey)
				if agentConn == nil {
					log.Printf("[Client] Agent not found: %s", connMsg.AgentKey)
					conn.SendSafe(&protocol.Message{
						Type: protocol.MessageTypeError,
					})
					continue
				}

				// Create session
				sessionID := uuid.New().String()
				s.hub.CreateSession(sessionID, conn.Connection, agentConn)

				// Send success
				conn.SendSafe(&protocol.Message{
					Type: protocol.MessageTypeConnect,
				})

				log.Printf("[Client] Session created: %s", sessionID)

			case protocol.MessageTypeHeartbeat:
				conn.SendSafe(&protocol.Message{Type: protocol.MessageTypePong})

			default:
				// Forward to hub (will route to agent via session)
				s.hub.messageClient <- &ClientMessage{
					ConnID: conn.ID,
					Msg:    &msg,
				}
			}
		}
	}

	// Cleanup
	close(done)
	s.hub.Unregister(conn.Connection)
	log.Printf("[Client] Connection closed: %s", conn.ID)
}
```

- [ ] **Step 4: Run tests to ensure everything compiles**

```bash
cd /root/projects/Weblise
go build ./server/...
go test ./server/ws/... -v
```

Expected: Compiles and tests pass

- [ ] **Step 5: Commit WebSocket server**

```bash
cd /root/projects/Weblise
git add server/ws/
git commit -m "feat: implement WebSocket server with agent and client handlers

- Add WebSocket server with /agent and /client endpoints
- Implement agent connection handler with register flow
- Implement client connection handler with connect flow
- Handle heartbeat/pong for keep-alive
- Route messages through hub between agents and clients

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 6: HTTP Server and Static Files

### Task 6: Implement HTTP Server with Static File Serving

**Files:**
- Create: `server/http/server.go`
- Create: `server/http/handler.go`
- Create: `server/http/static.go`

- [ ] **Step 1: Create HTTP server**

```go
// server/http/server.go
package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server manages the HTTP server
type Server struct {
	addr   string
	server *http.Server
}

// NewServer creates a new HTTP server
func NewServer(addr string) *Server {
	mux := http.NewServeMux()

	s := &Server{
		addr: addr,
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	// Register handlers
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/health", s.handleHealth)

	return s
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("[HTTP Server] Starting on %s", s.addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("[HTTP Server] Shutting down...")
	return s.server.Shutdown(ctx)
}

// handleIndex serves the main page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(indexHTML))
}

// handleHealth returns health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
```

- [ ] **Step 2: Create handler with WebSocket proxy**

```go
// server/http/handler.go
package http

import (
	"log"
	"net/http"
)

// handleWebSocket proxies WebSocket connections to the WebSocket server
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// For MVP, we'll have the WebSocket server on a separate port
	// In production, you'd use a reverse proxy or embed the WS server here
	http.Error(w, "WebSocket server on port 8443", http.StatusServiceUnavailable)
}

// handleAPI handles API endpoints
func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// API endpoints
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
```

- [ ] **Step 3: Create static files with embedded HTML/CSS/JS**

```go
// server/http/static.go
package http

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weblise - Remote Desktop</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #1a1a2e;
            color: #eee;
            min-height: 100vh;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }

        h1 {
            text-align: center;
            margin-bottom: 30px;
            color: #4ecca3;
        }

        #connection-panel {
            background: #16213e;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: flex;
            gap: 10px;
            align-items: center;
        }

        #agent-key {
            flex: 1;
            padding: 12px;
            border: 1px solid #0f3460;
            border-radius: 4px;
            background: #0f3460;
            color: #eee;
            font-size: 14px;
        }

        #agent-key:focus {
            outline: none;
            border-color: #4ecca3;
        }

        #connect-btn {
            padding: 12px 24px;
            background: #4ecca3;
            color: #1a1a2e;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-weight: bold;
            font-size: 14px;
        }

        #connect-btn:hover {
            background: #3db892;
        }

        #connect-btn:disabled {
            background: #555;
            cursor: not-allowed;
        }

        #status {
            padding: 10px;
            border-radius: 4px;
            margin-bottom: 20px;
            text-align: center;
            font-size: 14px;
        }

        .status-disconnected {
            background: #e94560;
        }

        .status-connecting {
            background: #f39c12;
        }

        .status-connected {
            background: #4ecca3;
            color: #1a1a2e;
        }

        #screen-container {
            display: none;
            background: #000;
            border-radius: 8px;
            overflow: hidden;
        }

        #screen-canvas {
            display: block;
            width: 100%;
            height: auto;
            cursor: crosshair;
        }

        .info {
            text-align: center;
            color: #666;
            margin-top: 20px;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Weblise Remote Desktop</h1>

        <div id="connection-panel">
            <input type="text" id="agent-key" placeholder="Enter Agent Key" autocomplete="off">
            <button id="connect-btn">Connect</button>
        </div>

        <div id="status" class="status-disconnected">Disconnected</div>

        <div id="screen-container">
            <canvas id="screen-canvas"></canvas>
        </div>

        <p class="info">Enter your agent key to connect to a remote device</p>
    </div>

    <script>
        const wsUrl = window.location.protocol === 'https:'
            ? 'wss://' + window.location.host + '/client'
            : 'ws://' + window.location.hostname + ':8443/client';

        class WebliseClient {
            constructor() {
                this.ws = null;
                this.connected = false;
                this.canvas = document.getElementById('screen-canvas');
                this.ctx = this.canvas.getContext('2d');

                this.bindEvents();
            }

            bindEvents() {
                document.getElementById('connect-btn').addEventListener('click', () => this.connect());
                document.getElementById('agent-key').addEventListener('keypress', (e) => {
                    if (e.key === 'Enter') this.connect();
                });

                this.canvas.addEventListener('mousemove', (e) => this.sendMouseMove(e));
                this.canvas.addEventListener('mousedown', (e) => this.sendMouseClick(e, 'down'));
                this.canvas.addEventListener('mouseup', (e) => this.sendMouseClick(e, 'up'));
            }

            connect() {
                const key = document.getElementById('agent-key').value.trim();
                if (!key) return;

                this.updateStatus('connecting', 'Connecting...');
                document.getElementById('connect-btn').disabled = true;

                try {
                    this.ws = new WebSocket(wsUrl);

                    this.ws.onopen = () => {
                        console.log('WebSocket connected');
                        // Send connect message
                        this.ws.send(JSON.stringify({
                            type: 'connect',
                            data: { agent_key: key }
                        }));
                    };

                    this.ws.onmessage = (event) => this.handleMessage(event);

                    this.ws.onclose = () => {
                        console.log('WebSocket disconnected');
                        this.connected = false;
                        this.updateStatus('disconnected', 'Disconnected');
                        document.getElementById('connect-btn').disabled = false;
                        document.getElementById('screen-container').style.display = 'none';
                    };

                    this.ws.onerror = (error) => {
                        console.error('WebSocket error:', error);
                        this.updateStatus('disconnected', 'Connection error');
                    };

                    // Start heartbeat
                    this.startHeartbeat();

                } catch (error) {
                    console.error('Connection failed:', error);
                    this.updateStatus('disconnected', 'Connection failed');
                    document.getElementById('connect-btn').disabled = false;
                }
            }

            handleMessage(event) {
                try {
                    const msg = JSON.parse(event.data);
                    console.log('Received:', msg.type);

                    switch (msg.type) {
                        case 'connect':
                            this.connected = true;
                            this.updateStatus('connected', 'Connected');
                            document.getElementById('screen-container').style.display = 'block';
                            break;

                        case 'frame':
                            this.renderFrame(msg.data);
                            break;

                        case 'error':
                            console.error('Server error:', msg.data?.message);
                            this.updateStatus('disconnected', 'Error: ' + (msg.data?.message || 'Unknown'));
                            break;

                        case 'pong':
                            // Heartbeat response
                            break;
                    }
                } catch (error) {
                    console.error('Failed to handle message:', error);
                }
            }

            renderFrame(data) {
                const img = new Image();
                img.onload = () => {
                    this.canvas.width = img.width;
                    this.canvas.height = img.height;
                    this.ctx.drawImage(img, 0, 0);
                };
                img.src = 'data:image/jpeg;base64,' + data;
            }

            sendMouseMove(e) {
                if (!this.connected) return;

                const rect = this.canvas.getBoundingClientRect();
                const scaleX = this.canvas.width / rect.width;
                const scaleY = this.canvas.height / rect.height;

                const x = Math.round((e.clientX - rect.left) * scaleX);
                const y = Math.round((e.clientY - rect.top) * scaleY);

                this.ws.send(JSON.stringify({
                    type: 'input',
                    data: {
                        action: 'mousemove',
                        data: { x, y }
                    }
                }));
            }

            sendMouseClick(e, action) {
                if (!this.connected) return;

                const rect = this.canvas.getBoundingClientRect();
                const scaleX = this.canvas.width / rect.width;
                const scaleY = this.canvas.height / rect.height;

                const x = Math.round((e.clientX - rect.left) * scaleX);
                const y = Math.round((e.clientY - rect.top) * scaleY);
                const button = e.button;

                this.ws.send(JSON.stringify({
                    type: 'input',
                    data: {
                        action: 'mouse' + action,
                        data: { x, y, button }
                    }
                }));
            }

            updateStatus(state, message) {
                const statusEl = document.getElementById('status');
                statusEl.textContent = message;
                statusEl.className = 'status-' + state;
            }

            startHeartbeat() {
                setInterval(() => {
                    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                        this.ws.send(JSON.stringify({ type: 'ping' }));
                    }
                }, 30000);
            }
        }

        // Initialize client
        const client = new WebliseClient();
    </script>
</body>
</html>
`
```

- [ ] **Step 4: Test HTTP server compiles**

```bash
cd /root/projects/Weblise
go build ./server/http/...
```

Expected: Compiles without errors

- [ ] **Step 5: Commit HTTP server**

```bash
cd /root/projects/Weblise
git add server/http/
git commit -m "feat: implement HTTP server with embedded web client

- Add HTTP server for serving static files
- Embed complete web client in a single HTML file
- Implement WebSocket connection handling in JavaScript
- Add mouse input handling
- Add responsive UI with connection status

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 7: Server Main Entry Point

### Task 7: Implement Server Main Entry Point

**Files:**
- Modify: `server/main.go`
- Create: `server/cmd/root.go`
- Create: `server/config/config.go`

- [ ] **Step 1: Create config**

```go
// server/config/config.go
package config

import (
	"fmt"
	"os"
)

// Config holds server configuration
type Config struct {
	HTTPPort string
	WSPort   string
	DBPath   string
}

// Load loads configuration from environment or defaults
func Load() *Config {
	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		WSPort:   getEnv("WS_PORT", "8443"),
		DBPath:   getEnv("DB_PATH", "./data/weblise.db"),
	}
}

// HTTPAddr returns the HTTP server address
func (c *Config) HTTPAddr() string {
	return fmt.Sprintf(":%s", c.HTTPPort)
}

// WSAddr returns the WebSocket server address
func (c *Config) WSAddr() string {
	return fmt.Sprintf(":%s", c.WSPort)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

- [ ] **Step 2: Update server main.go**

```go
// server/main.go
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

	// Load config
	cfg := config.Load()
	log.Printf("Config: HTTP=%s, WS=%s, DB=%s", cfg.HTTPPort, cfg.WSPort, cfg.DBPath)

	// Ensure data directory
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize database
	database, err := db.New(db.Config{Path: cfg.DBPath})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("Database initialized")

	// Create hub
	hub := ws.NewHub()
	go hub.Run()

	log.Println("Hub started")

	// Start WebSocket server
	wsServer := ws.NewServer(cfg.WSAddr(), hub)
	go func() {
		if err := wsServer.Start(); err != nil {
			log.Fatalf("WebSocket server failed: %v", err)
		}
	}()

	// Start HTTP server
	httpServer := http.NewServer(cfg.HTTPAddr())
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Println("Weblise Server is running")
	log.Printf("  - HTTP:  http://localhost%s", cfg.HTTPAddr())
	log.Printf("  - WebSocket: ws://localhost%s", cfg.WSAddr())

	// Wait for interrupt signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")

	// Graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// hub.Stop()

	log.Println("Server stopped")
}
```

- [ ] **Step 3: Test server compiles**

```bash
cd /root/projects/Weblise/server
go build -o server main.go
```

Expected: Creates server binary

- [ ] **Step 4: Test server starts**

```bash
cd /root/projects/Weblise/server
timeout 3 ./server || true
```

Expected output:
```
Weblise Server v0.1.0
Config: HTTP=8080, WS=8443, DB=./data/weblise.db
Database initialized
Hub started
[WS Server] Starting on :8443
[HTTP Server] Starting on :8080
Weblise Server is running
```

- [ ] **Step 5: Commit server main**

```bash
cd /root/projects/Weblise
git add server/
git commit -m "feat: implement server main entry point

- Add config loading from environment variables
- Initialize database and hub
- Start HTTP and WebSocket servers
- Add graceful shutdown handling
- Add data directory creation

Co-Authored-By: Claude <noreply@anthropic.com>"
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 8: Agent Screen Capture

### Task 8: Implement Agent Screen Capture

**Files:**
- Create: `agent/screen/capture.go`
- Create: `agent/screen/capture_windows.go`
- Create: `agent/screen/capture_linux.go`

- [ ] **Step 1: Create screen capture interface**

```go
// agent/screen/capture.go
package screen

import (
	"image"
	"time"
)

// Capturer defines the interface for screen capture
type Capturer interface {
	Capture() (*image.RGBA, error)
	Bounds() (width, height int)
}

// Config holds screen capture configuration
type Config struct {
	Quality int    // JPEG quality (1-100)
	Format  string // "jpeg" or "png"
}

// DefaultConfig returns default capture config
func DefaultConfig() *Config {
	return &Config{
		Quality: 80,
		Format:  "jpeg",
	}
}

// New creates a new screen capturer for the current platform
func New(cfg *Config) (Capturer, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return newCapturer(cfg)
}

// Frame represents a captured frame
type Frame struct {
	Data      []byte
	Width     int
	Height    int
	Timestamp time.Time
}
```

- [ ] **Step 2: Create Windows screen capture**

```go
// agent/screen/capture_windows.go
// +build windows

package screen

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"runtime"
	"time"

	"github.com/kbinani/screenshot"
)

type windowsCapturer struct {
	config     *Config
	screenNum  int
	lastBounds image.Rectangle
}

func newCapturer(cfg *Config) (Capturer, error) {
	c := &windowsCapturer{
		config:    cfg,
		screenNum: 0, // Primary screen
	}

	// Validate screen exists
	n := screenshot.NumActiveDisplays()
	if n == 0 {
		return nil, os.ErrNotExist
	}

	bounds := screenshot.GetDisplayBounds(c.screenNum)
	c.lastBounds = bounds

	return c, nil
}

func (c *windowsCapturer) Capture() (*image.RGBA, error) {
	bounds := screenshot.GetDisplayBounds(c.screenNum)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	c.lastBounds = bounds
	return img, nil
}

func (c *windowsCapturer) Bounds() (int, int) {
	bounds := screenshot.GetDisplayBounds(c.screenNum)
	return bounds.Dx(), bounds.Dy()
}

// EncodeToBytes encodes the image to bytes
func EncodeFrame(img *image.RGBA, format string, quality int) ([]byte, error) {
	if format == "png" {
		// For MVP, use PNG for simplicity
		// Later: optimize with JPEG
		return encodePNG(img)
	}
	return encodeJPEG(img, quality)
}

func encodeJPEG(img *image.RGBA, quality int) ([]byte, error) {
	// JPEG encoding not directly supported for image.RGBA
	// Use PNG for MVP
	return encodePNG(img)
}

func encodePNG(img *image.RGBA) ([]byte, error) {
	// Write to buffer
	// For simplicity, return base64 encoded data
	// The actual encoding will be done in the connection layer
	return nil, nil
}
```

- [ ] **Step 3: Create Linux screen capture**

```go
// agent/screen/capture_linux.go
// +build linux

package screen

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type linuxCapturer struct {
	config *Config
}

func newCapturer(cfg *Config) (Capturer, error) {
	// Check if we're on X11
	if os.Getenv("DISPLAY") == "" {
		return nil, fmt.Errorf("DISPLAY environment variable not set")
	}

	// Check for required tools
	if _, err := exec.LookPath("scrot"); err != nil {
		if _, err2 := exec.LookPath("import"); err2 != nil {
			return nil, fmt.Errorf("no screenshot tool found (need scrot or imagemagick)")
		}
	}

	return &linuxCapturer{config: cfg}, nil
}

func (c *linuxCapturer) Capture() (*image.RGBA, error) {
	// For MVP, use scrot if available
	// This is a placeholder - actual implementation would capture screen
	return nil, fmt.Errorf("not implemented yet - use scrot or x11 API")
}

func (c *linuxCapturer) Bounds() (int, int) {
	// Return default resolution
	return 1920, 1080
}
```

- [ ] **Step 4: Add simple JPEG encoder**

```go
// Add to agent/screen/capture.go

import (
	"bytes"
	"image/jpeg"
	"image/png"
)

// EncodeJPEG encodes image to JPEG bytes
func EncodeJPEG(img *image.RGBA, quality int) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodePNG encodes image to PNG bytes
func EncodePNG(img *image.RGBA) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeToBase64 encodes bytes to base64
import "encoding/base64"

func EncodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
```

- [ ] **Step 5: Create screen capture test**

```go
// agent/screen/capture_test.go
package screen

import (
	"testing"
)

func TestNewCapturer(t *testing.T) {
	cfg := DefaultConfig()
	capturer, err := New(cfg)
	if err != nil {
		t.Logf("Capturer creation failed (expected on some systems): %v", err)
		return
	}
	defer capturer.Close()

	width, height := capturer.Bounds()
	if width <= 0 || height <= 0 {
		t.Errorf("Invalid bounds: %dx%d", width, height)
	}

	t.Logf("Screen bounds: %dx%d", width, height)
}

func TestCapture(t *testing.T) {
	cfg := DefaultConfig()
	capturer, err := New(cfg)
	if err != nil {
		t.Skip("Capturer not available:", err)
		return
	}
	defer capturer.Close()

	img, err := capturer.Capture()
	if err != nil {
		t.Fatalf("Capture failed: %v", err)
	}

	if img == nil {
		t.Fatal("Image is nil")
	}

	if img.Bounds().Dx() <= 0 || img.Bounds().Dy() <= 0 {
		t.Error("Invalid image size")
	}

	// Test encoding
	data, err := EncodeJPEG(img, 80)
	if err != nil {
		t.Fatalf("JPEG encoding failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Encoded data is empty")
	}

	t.Logf("Captured and encoded %d bytes", len(data))
}
```

- [ ] **Step 6: Fix build tags and test**

```bash
cd /root/projects/Weblise
go test ./agent/screen/... -v
```

- [ ] **Step 7: Commit screen capture**

```bash
cd /root/projects/Weblise
git add agent/screen/
git commit -m "feat: implement screen capture for agent

- Add screen capture interface
- Implement Windows capture using screenshot library
- Implement Linux capture stub
- Add JPEG/PNG encoding utilities
- Add unit tests

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 9: Agent Input Control

### Task 9: Implement Agent Input Control

**Files:**
- Create: `agent/input/mouse.go`
- Create: `agent/input/keyboard.go`

- [ ] **Step 1: Create mouse control**

```go
// agent/input/mouse.go
package input

import (
	"github.com/go-vgo/robotgo"
)

// Mouse handles mouse input
type Mouse struct{}

// NewMouse creates a new mouse controller
func NewMouse() *Mouse {
	return &Mouse{}
}

// Move moves the mouse to absolute position
func (m *Mouse) Move(x, y int) error {
	return robotgo.MoveMouse(x, y)
}

// Click clicks a mouse button
func (m *Mouse) Click(button string) error {
	var buttonStr string
	switch button {
	case "left":
		buttonStr = "left"
	case "right":
		buttonStr = "right"
	case "middle":
		buttonStr = "center"
	default:
		buttonStr = "left"
	}
	return robotgo.MouseClick(buttonStr)
}

// Press presses a mouse button down
func (m *Mouse) Press(button string) error {
	var mouseButton string
	switch button {
	case "left", "0":
		mouseButton = "left"
	case "right", "2":
		mouseButton = "right"
	case "middle", "1":
		mouseButton = "center"
	default:
		mouseButton = "left"
	}
	return robotgo.MouseToggle("down", mouseButton)
}

// Release releases a mouse button
func (m *Mouse) Release(button string) error {
	var mouseButton string
	switch button {
	case "left", "0":
		mouseButton = "left"
	case "right", "2":
		mouseButton = "right"
	case "middle", "1":
		mouseButton = "center"
	default:
		mouseButton = "left"
	}
	return robotgo.MouseToggle("up", mouseButton)
}

// Scroll scrolls the mouse wheel
func (m *Mouse) Scroll(x, y int) error {
	return robotgo.Scroll(x, y)
}
```

- [ ] **Step 2: Create keyboard control**

```go
// agent/input/keyboard.go
package input

import (
	"github.com/go-vgo/robotgo"
)

// Keyboard handles keyboard input
type Keyboard struct{}

// NewKeyboard creates a new keyboard controller
func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

// Press presses a key down
func (k *Keyboard) Press(key string) error {
	return robotgo.KeyTap(key)
}

// Type types a string
func (k *Keyboard) Type(text string) error {
	return robotgo.TypeStr(text)
}

// Release releases a key
func (k *Keyboard) Release(key string) error {
	// robotgo doesn't have a direct release, Tap handles press+release
	return nil
}
```

- [ ] **Step 3: Create input tests**

```go
// agent/input/mouse_test.go
package input

import (
	"testing"
)

func TestMouseMove(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping mouse test in short mode")
	}

	mouse := NewMouse()
	if err := mouse.Move(100, 100); err != nil {
		t.Errorf("MouseMove failed: %v", err)
	}
}

func TestMouseClick(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping mouse test in short mode")
	}

	mouse := NewMouse()
	if err := mouse.Click("left"); err != nil {
		t.Errorf("MouseClick failed: %v", err)
	}
}
```

- [ ] **Step 4: Commit input control**

```bash
cd /root/projects/Weblise
git add agent/input/
git commit -m "feat: implement mouse and keyboard input control

- Add mouse controller using robotgo
- Implement move, click, press, release
- Add keyboard controller for typing
- Add unit tests

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 10: Agent WebSocket Connection

### Task 10: Implement Agent WebSocket Connection

**Files:**
- Create: `agent/conn/ws.go`
- Modify: `agent/main.go`

- [ ] **Step 1: Create agent WebSocket connection**

```go
// agent/conn/ws.go
package conn

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/chenjingxiong/weblise/server/protocol"
	"github.com/chenjingxiong/weblise/agent/input"
	"github.com/chenjingxiong/weblise/agent/screen"
	"github.com/gorilla/websocket"
)

// Connection manages the WebSocket connection to the server
type Connection struct {
	serverAddr string
	deviceKey  string
	name       string
	wsConn     *websocket.Conn
	capturer   screen.Capturer
	mouse      *input.Mouse
	keyboard   *input.Keyboard

	mu             sync.Mutex
	sendChan       chan *protocol.Message
	connected      bool
	screenInterval time.Duration
}

// Config holds connection configuration
type Config struct {
	ServerAddr     string
	DeviceKey      string
	Name           string
	CaptureFPS     int
	ScreenQuality  int
}

// New creates a new connection
func New(cfg *Config) (*Connection, error) {
	// Create screen capturer
	capturer, err := screen.New(screen.DefaultConfig())
	if err != nil {
		return nil, err
	}

	// Calculate capture interval
	interval := time.Second / time.Duration(cfg.CaptureFPS)

	return &Connection{
		serverAddr:     cfg.ServerAddr,
		deviceKey:      cfg.DeviceKey,
		name:          cfg.Name,
		capturer:      capturer,
		mouse:         input.NewMouse(),
		keyboard:      input.NewKeyboard(),
		sendChan:      make(chan *protocol.Message, 32),
		screenInterval: interval,
	}, nil
}

// Connect establishes the WebSocket connection
func (c *Connection) Connect() error {
	wsConn, _, err := websocket.DefaultDialer.Dial(c.serverAddr, nil)
	if err != nil {
		return err
	}

	c.wsConn = wsConn
	c.connected = true

	log.Printf("[Conn] Connected to %s", c.serverAddr)

	// Send register message
	regMsg, _ := protocol.NewMessage(protocol.MessageTypeRegister, protocol.RegisterMessage{
		Key:   c.deviceKey,
		Name:  c.name,
		OSType: getOSType(),
	})

	if err := c.wsConn.WriteJSON(regMsg); err != nil {
		return err
	}

	log.Printf("[Conn] Registered with key: %s", c.deviceKey)

	return nil
}

// Start starts the connection's main loops
func (c *Connection) Start() error {
	// Start send loop
	go c.sendLoop()

	// Start receive loop
	go c.receiveLoop()

	// Start screen capture loop
	go c.screenLoop()

	// Start heartbeat
	go c.heartbeatLoop()

	return nil
}

// sendLoop handles sending messages
func (c *Connection) sendLoop() {
	for {
		select {
		case msg := <-c.sendChan:
			if !c.connected {
				return
			}
			if err := c.wsConn.WriteJSON(msg); err != nil {
				log.Printf("[Conn] Send error: %v", err)
				c.Close()
				return
			}
		}
	}
}

// receiveLoop handles incoming messages
func (c *Connection) receiveLoop() {
	for {
		if !c.connected {
			return
		}

		var msg protocol.Message
		if err := c.wsConn.ReadJSON(&msg); err != nil {
			log.Printf("[Conn] Receive error: %v", err)
			c.Close()
			return
		}

		c.handleMessage(&msg)
	}
}

// handleMessage handles an incoming message
func (c *Connection) handleMessage(msg *protocol.Message) {
	switch msg.Type {
	case protocol.MessageTypeInput:
		c.handleInput(msg)
	case protocol.MessageTypePong:
		// Heartbeat response
	case protocol.MessageTypeError:
		log.Printf("[Conn] Server error: %v", msg)
	default:
		log.Printf("[Conn] Unknown message type: %s", msg.Type)
	}
}

// handleInput handles an input message
func (c *Connection) handleInput(msg *protocol.Message) {
	var inputMsg protocol.InputMessage
	if err := msg.ParseData(&inputMsg); err != nil {
		log.Printf("[Conn] Parse input failed: %v", err)
		return
	}

	switch inputMsg.Action {
	case "mousemove":
		c.mouse.Move(inputMsg.Data.X, inputMsg.Data.Y)

	case "mousedown", "mouseup":
		button := getButtonName(inputMsg.Data.Button)
		if inputMsg.Action == "mousedown" {
			c.mouse.Press(button)
		} else {
			c.mouse.Release(button)
		}

	case "keydown":
		c.keyboard.Press(inputMsg.Data.Key)

	case "keyup":
		c.keyboard.Release(inputMsg.Data.Key)
	}
}

// screenLoop captures and sends screen frames
func (c *Connection) screenLoop() {
	ticker := time.NewTicker(c.screenInterval)
	defer ticker.Stop()

	for range ticker.C {
		if !c.connected {
			return
		}

		// Capture screen
		img, err := c.capturer.Capture()
		if err != nil {
			log.Printf("[Conn] Capture failed: %v", err)
			continue
		}

		// Encode to JPEG
		data, err := screen.EncodeJPEG(img, 80)
		if err != nil {
			log.Printf("[Conn] Encode failed: %v", err)
			continue
		}

		// Encode to base64
		base64Data := base64.StdEncoding.EncodeToString(data)

		// Send frame message
		frameMsg := protocol.FrameMessage{
			Data:      []byte(base64Data),
			Width:     img.Bounds().Dx(),
			Height:    img.Bounds().Dy(),
			Timestamp: time.Now().Unix(),
		}

		msg, _ := protocol.NewMessage(protocol.MessageTypeFrame, frameMsg)
		c.sendChan <- msg
	}
}

// heartbeatLoop sends periodic heartbeats
func (c *Connection) heartbeatLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !c.connected {
			return
		}

		pingMsg := protocol.HeartbeatMessage{
			Timestamp: time.Now().Unix(),
		}

		msg, _ := protocol.NewMessage(protocol.MessageTypeHeartbeat, pingMsg)
		c.sendChan <- msg
	}
}

// Close closes the connection
func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return
	}

	c.connected = false
	if c.wsConn != nil {
		c.wsConn.Close()
	}
	close(c.sendChan)

	log.Printf("[Conn] Connection closed")
}

// getButtonName converts button number to name
func getButtonName(button int) string {
	switch button {
	case 0:
		return "left"
	case 1:
		return "middle"
	case 2:
		return "right"
	default:
		return "left"
	}
}

// getOSType returns the current OS type
func getOSType() string {
	// Simple detection
	return "linux" // Will be platform-specific in build
}
```

- [ ] **Step 2: Update agent main.go**

```go
// agent/main.go
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

	// Generate default name if not provided
	deviceName := *name
	if deviceName == "" {
		hostname, _ := os.Hostname()
		deviceName = hostname
	}

	// Create connection
	connection, err := conn.New(&conn.Config{
		ServerAddr:    *serverAddr,
		DeviceKey:     *deviceKey,
		Name:         deviceName,
		CaptureFPS:   *fps,
		ScreenQuality: *quality,
	})
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}

	// Connect to server
	if err := connection.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Start connection loops
	connection.Start()

	log.Println("Agent is running, press Ctrl+C to exit")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	connection.Close()
	log.Println("Agent stopped")
}
```

- [ ] **Step 3: Fix imports in connection**

```go
// Add import for screen package functions
// The screen package should export EncodeJPEG
```

- [ ] **Step 4: Test agent compiles**

```bash
cd /root/projects/Weblise/agent
go build -o agent main.go
```

Expected: Compiles without errors

- [ ] **Step 5: Commit agent connection**

```bash
cd /root/projects/Weblise
git add agent/
git commit -m "feat: implement agent WebSocket connection

- Add WebSocket connection to server
- Implement register flow with device key
- Add screen capture loop with configurable FPS
- Add input message handling for mouse/keyboard
- Add heartbeat for connection keep-alive
- Update main.go with CLI flags

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 11: Docker Deployment

### Task 11: Create Docker Configuration

**Files:**
- Create: `docker/Dockerfile`
- Create: `docker/docker-compose.yml`

- [ ] **Step 1: Create Dockerfile**

```dockerfile
# docker/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build server
RUN CGO_ENABLED=1 go build -o weblise-server ./server

# Final image
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Copy binary from builder
COPY --from=builder /build/weblise-server /app/weblise-server

# Create data directory
RUN mkdir -p /app/data

# Set environment variables
ENV HTTP_PORT=8080
ENV WS_PORT=8443
ENV DB_PATH=/app/data/weblise.db

# Expose ports
EXPOSE 8080 8443

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O- http://localhost:8080/health || exit 1

# Run the server
CMD ["./weblise-server"]
```

- [ ] **Step 2: Create docker-compose.yml**

```yaml
# docker/docker-compose.yml
version: '3.8'

services:
  weblise-server:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    container_name: weblise-server
    ports:
      - "8080:8080"   # HTTP for reverse proxy
      - "8443:8443"   # WebSocket for agents
    volumes:
      - ./data:/app/data  # SQLite data persistence
    environment:
      - HTTP_PORT=8080
      - WS_PORT=8443
      - DB_PATH=/app/data/weblise.db
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "-O-", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

- [ ] **Step 3: Add .dockerignore**

```dockerignore
# .dockerignore
*.md
.git
.gitignore
data/
agent/
web/
docs/
scripts/
```

- [ ] **Step 4: Create build script**

```bash
#!/bin/bash
# scripts/build-docker.sh

set -e

echo "Building Weblise Docker image..."
docker build -f docker/Dockerfile -t weblise-server:latest .

echo "Done! Image built: weblise-server:latest"
echo ""
echo "To run:"
echo "  docker run -p 8080:8080 -p 8443:8443 -v \$(pwd)/data:/app/data weblise-server:latest"
echo ""
echo "Or use docker-compose:"
echo "  cd docker && docker-compose up -d"
```

- [ ] **Step 5: Make build script executable**

```bash
chmod +x /root/projects/Weblise/scripts/build-docker.sh
```

- [ ] **Step 6: Create deployment documentation**

Add to README.md:

```markdown
## Deployment

### Docker

Build the image:

```bash
./scripts/build-docker.sh
```

Run with Docker:

```bash
docker run -d \
  -p 8080:8080 \
  -p 8443:8443 \
  -v $(pwd)/data:/app/data \
  -e HTTP_PORT=8080 \
  -e WS_PORT=8443 \
  --name weblise \
  weblise-server:latest
```

Run with Docker Compose:

```bash
cd docker
docker-compose up -d
```

### Reverse Proxy (Nginx)

Example Nginx configuration:

\`\`\`nginx
server {
    listen 443 ssl;
    server_name remote.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # Web client
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # WebSocket
    location /ws/ {
        proxy_pass http://localhost:8080/ws/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # Agent WebSocket (separate port)
    location /agent/ {
        proxy_pass http://localhost:8443/agent/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
\`\`\`
```

- [ ] **Step 7: Commit Docker configuration**

```bash
cd /root/projects/Weblise
git add docker/ .dockerignore scripts/ README.md
git commit -m "feat: add Docker deployment configuration

- Add Dockerfile for server
- Add docker-compose.yml for easy deployment
- Add .dockerignore to exclude unnecessary files
- Add build script for Docker image
- Add deployment documentation to README

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 12: End-to-End Testing

### Task 12: Create End-to-End Tests

**Files:**
- Create: `e2e/main_test.go`

- [ ] **Step 1: Create E2E test framework**

```go
// e2e/main_test.go
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

const (
	serverAddr = "http://localhost:8080"
	wsAddr     = "ws://localhost:8443"
)

func TestMain(m *testing.M) {
	// Start server
	if err := startServer(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}

	// Wait for server to be ready
	if !waitForServer(serverAddr + "/health") {
		fmt.Println("Server failed to start")
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	stopServer()

	os.Exit(code)
}

func startServer() error {
	cmd := exec.Command("go", "run", "../server/main.go")
	cmd.Dir = ".."
	if err := cmd.Start(); err != nil {
		return err
	}
	serverCmd = cmd
	return nil
}

var serverCmd *exec.Cmd

func stopServer() {
	if serverCmd != nil {
		serverCmd.Process.Kill()
	}
}

func waitForServer(url string) bool {
	for i := 0; i < 30; i++ {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return resp.StatusCode == http.StatusOK
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get(serverAddr + "/health")
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result["status"])
	}
}

func TestIndexPage(t *testing.T) {
	resp, err := http.Get(serverAddr + "/")
	if err != nil {
		t.Fatalf("Index request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if len(body) == 0 {
		t.Error("Response body is empty")
	}

	content := string(body)
	if !contains(content, "Weblise") {
		t.Error("Response doesn't contain 'Weblise'")
	}
	if !contains(content, "agent-key") {
		t.Error("Response doesn't contain agent-key input")
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
```

- [ ] **Step 2: Create WebSocket E2E test**

```go
// e2e/websocket_test.go
package e2e

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestAgentConnection(t *testing.T) {
	// Connect to agent endpoint
	ws, _, err := websocket.DefaultDialer.Dial(wsAddr+"/agent", nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer ws.Close()

	// Send register message
	regMsg := map[string]interface{}{
		"type": "register",
		"data": map[string]string{
			"key":     "test-agent-e2e",
			"name":    "E2E Test Agent",
			"os_type": "linux",
		},
	}

	if err := ws.WriteJSON(regMsg); err != nil {
		t.Fatalf("Write register failed: %v", err)
	}

	// Read response
	var msg map[string]interface{}
	if err := ws.ReadJSON(&msg); err != nil {
		t.Fatalf("Read response failed: %v", err)
	}

	// Check if connected successfully
	if msg["type"] == nil {
		t.Error("No message type in response")
	}
}

func TestClientConnection(t *testing.T) {
	// First, register an agent
	agentWS, _, err := websocket.DefaultDialer.Dial(wsAddr+"/agent", nil)
	if err != nil {
		t.Skip("Agent connection failed, skipping client test")
		return
	}

	regMsg := map[string]interface{}{
		"type": "register",
		"data": map[string]string{
			"key":  "test-agent-client-e2e",
			"name": "E2E Test Agent",
		},
	}
	agentWS.WriteJSON(regMsg)

	// Wait a bit for registration
	time.Sleep(100 * time.Millisecond)

	// Connect client
	clientWS, _, err := websocket.DefaultDialer.Dial(wsAddr+"/client", nil)
	if err != nil {
		t.Fatalf("Client dial failed: %v", err)
	}
	defer clientWS.Close()

	// Send connect message
	connMsg := map[string]interface{}{
		"type": "connect",
		"data": map[string]string{
			"agent_key": "test-agent-client-e2e",
		},
	}

	if err := clientWS.WriteJSON(connMsg); err != nil {
		t.Fatalf("Write connect failed: %v", err)
	}

	// Read response
	var msg map[string]interface{}
	if err := clientWS.ReadJSON(&msg); err != nil {
		t.Fatalf("Read response failed: %v", err)
	}

	t.Logf("Client received: %v", msg)

	agentWS.Close()
}
```

- [ ] **Step 3: Run E2E tests**

```bash
cd /root/projects/Weblise
go test ./e2e/... -v -timeout 30s
```

Expected: Tests pass

- [ ] **Step 4: Commit E2E tests**

```bash
cd /root/projects/Weblise
git add e2e/
git commit -m "test: add end-to-end tests

- Add E2E test framework with server startup/shutdown
- Test health endpoint
- Test index page loads correctly
- Test agent WebSocket connection
- Test client WebSocket connection

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Chunk 13: Documentation and Finalization

### Task 13: Update Documentation

**Files:**
- Modify: `README.md`
- Create: `docs/DEPLOYMENT.md`

- [ ] **Step 1: Update README.md**

```markdown
# Weblise

Web-based remote desktop system - Connect to any device from your browser.

![Version](https://img.shields.io/badge/version-0.1.0-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)

## Features

- **Web-based**: Control devices from any modern browser
- **Simple deployment**: Single Docker container
- **Secure**: Works behind HTTPS reverse proxy
- **Cross-platform agent**: Windows and Linux support
- **Low bandwidth**: Optimized JPEG compression

## Quick Start

### 1. Start the Server

\`\`\`bash
docker run -d \\
  -p 8080:8080 \\
  -p 8443:8443 \\
  -v \$(pwd)/data:/app/data \\
  ghcr.io/chenjingxiong/weblise-server:latest
\`\`\`

### 2. Connect an Agent

\`\`\`bash
# Download agent
wget https://github.com/chenjingxiong/weblise/releases/latest/download/agent-linux-amd64
chmod +x agent-linux-amd64

# Start agent
./agent-linux-amd64 --server=ws://your-server.com:8443 --key=your-device-key
\`\`\`

### 3. Control from Browser

Open https://your-server.com and enter your device key.

## Architecture

```
┌─────────┐                    ┌─────────┐
│ Client  │                    │  Agent  │
│(Browser)│◄───WebSocket──────►│ (Go)    │
└────┬────┘                    └────┬────┘
     │                              │
     └──────────┬───────────────────┘
                │
         ┌──────▼──────┐
         │   Server    │
         │    (Go)     │
         │  + SQLite   │
         └─────────────┘
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | 8080 | HTTP server port |
| `WS_PORT` | 8443 | WebSocket server port |
| `DB_PATH` | ./data/weblise.db | SQLite database path |

### Agent Flags

| Flag | Description |
|------|-------------|
| `--server` | WebSocket server address (required) |
| `--key` | Device authentication key (required) |
| `--name` | Device name (default: hostname) |
| `--fps` | Screen capture FPS (default: 15) |
| `--quality` | JPEG quality 1-100 (default: 80) |

## Development

See [DEVELOPMENT.md](docs/DEVELOPMENT.md) for development setup.

## Deployment

See [DEPLOYMENT.md](docs/DEPLOYMENT.md) for production deployment guide.

## License

MIT

## Roadmap

- [ ] P2P connections for lower latency
- [ ] User authentication system
- [ ] Mobile apps (iOS/Android)
- [ ] File transfer
- [ ] Clipboard sync
- [ ] Multi-monitor support
```

- [ ] **Step 2: Create deployment guide**

```markdown
# docs/DEPLOYMENT.md

# Deployment Guide

## Production Deployment

### Prerequisites

- Docker and Docker Compose
- Domain name with SSL certificate
- Reverse proxy (Nginx, Caddy, Traefik)

### 1. Deploy Server

\`\`\`bash
git clone https://github.com/chenjingxiong/weblise.git
cd weblise/docker
docker-compose up -d
\`\`\`

### 2. Configure Reverse Proxy

#### Nginx

\`\`\`nginx
server {
    listen 443 ssl http2;
    server_name remote.example.com;

    ssl_certificate /etc/ssl/certs/example.com.crt;
    ssl_certificate_key /etc/ssl/private/example.com.key;

    # Web client
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket for clients
    location /ws/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # WebSocket for agents
    location /agent/ {
        proxy_pass http://localhost:8443;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
\`\`\`

#### Caddy

\`\`\`caddy
remote.example.com {
    reverse_proxy localhost:8080

    handle /ws/* {
        reverse_proxy localhost:8080
    }

    handle /agent/* {
        reverse_proxy localhost:8443
    }
}
\`\`\`

### 3. Deploy Agents

\`\`\`bash
# Download agent for Linux
wget https://github.com/chenjingxiong/weblise/releases/latest/download/agent-linux-amd64
sudo mv agent-linux-amd64 /usr/local/bin/weblise-agent
sudo chmod +x /usr/local/bin/weblise-agent

# Create systemd service
sudo tee /etc/systemd/system/weblise-agent.service > /dev/null <<EOF
[Unit]
Description=Weblise Agent
After=network.target

[Service]
Type=simple
User=nobody
ExecStart=/usr/local/bin/weblise-agent \\
    --server=ws://remote.example.com:8443 \\
    --key=YOUR_DEVICE_KEY \\
    --name=\$(hostname)
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
sudo systemctl enable weblise-agent
sudo systemctl start weblise-agent
\`\`\`

### 4. Firewall

Ensure ports are accessible:

\`\`\`bash
# If direct access (not using reverse proxy for agent port)
sudo ufw allow 8443/tcp
\`\`\`

## Monitoring

### Health Check

\`\`\`bash
curl http://localhost:8080/health
\`\`\`

### Logs

\`\`\`bash
docker logs -f weblise-server
\`\`\`

## Security

- Always use HTTPS in production
- Use strong, unique device keys
- Consider rate limiting for connection attempts
- Keep agent keys secret
- Regular security updates
```

- [ ] **Step 3: Create CHANGELOG.md**

```markdown
# CHANGELOG.md

## [0.1.0] - 2026-03-30

### Added
- Initial MVP release
- Agent for Windows/Linux with screen capture
- Server with WebSocket routing
- Web client for browser-based control
- Mouse and keyboard input
- Docker deployment
- SQLite for device management
```

- [ ] **Step 4: Create development docs**

```markdown
# docs/DEVELOPMENT.md

# Development Guide

## Setup

\`\`\`bash
# Clone repository
git clone https://github.com/chenjingxiong/weblise.git
cd weblise

# Install dependencies
go mod download

# Run tests
go test ./...

# Run server
cd server
go run main.go
\`\`\`

## Project Structure

- `agent/` - Agent code (screen capture, input control)
- `server/` - Server code (WebSocket routing, database)
- `web/` - Web client files
- `docker/` - Docker configuration
- `e2e/` - End-to-end tests

## Running Locally

### Server

\`\`\`bash
cd server
go run main.go
\`\`\`

### Agent

\`\`\`bash
cd agent
go run main.go --server=ws://localhost:8443 --key=test-key
\`\`\`

### Testing

\`\`\`bash
# Unit tests
go test ./...

# E2E tests
go test ./e2e/... -v
\`\`\`

## Building

\`\`\`bash
# Server
cd server
go build -o weblise-server main.go

# Agent
cd agent
GOOS=linux GOARCH=amd64 go build -o agent-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o agent-windows-amd64.exe main.go
\`\`\`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request
```

- [ ] **Step 5: Commit documentation**

```bash
cd /root/projects/Weblise
git add README.md docs/ CHANGELOG.md
git commit -m "docs: update project documentation

- Update README with quick start and features
- Add DEPLOYMENT.md with production guide
- Add DEVELOPMENT.md with development setup
- Add CHANGELOG.md for version tracking
- Update README with architecture diagram

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

---

## Completion Checklist

- [ ] All tests pass: `go test ./...`
- [ ] Server builds: `cd server && go build`
- [ ] Agent builds: `cd agent && go build`
- [ ] Docker image builds: `docker build -f docker/Dockerfile`
- [ ] E2E tests pass: `go test ./e2e/...`
- [ ] Documentation complete
- [ ] README has quick start
- [ ] Deployment guide exists

---

**Next Steps After MVP:**

1. Add P2P connection support with WebRTC data channels
2. Implement user authentication system
3. Add native mobile clients (Flutter)
4. Optimize screen capture with hardware encoding
5. Add file transfer capabilities
6. Implement multi-monitor support
7. Add clipboard synchronization

---

**Implementation Notes:**

- This is an MVP focused on core functionality
- Security can be enhanced in future versions (E2E encryption, rate limiting)
- Performance optimizations can be added as needed
- The protocol is extensible for future features
