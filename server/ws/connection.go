package ws

import (
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
