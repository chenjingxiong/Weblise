package ws

import (
	"log"
	"sync"

	"github.com/chenjingxiong/weblise/server/protocol"
)

// Hub manages active connections and routes messages between agents and clients
type Hub struct {
	mu            sync.RWMutex
	agents        map[string]*Connection // key: connection ID
	clients       map[string]*Connection // key: connection ID
	agentsByKey   map[string]*Connection // key: device key
	sessions      map[string]*Session    // key: session ID
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
					session.ClientConn.SendSafe(&protocol.Message{Type: protocol.MessageTypeError})
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
