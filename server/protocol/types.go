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
	ID         string
	Type       ConnectionType
	DeviceKey  string
	ConnectedAt time.Time
	LastSeen   time.Time
	SendChan   chan *Message
	CloseChan  chan struct{}
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
