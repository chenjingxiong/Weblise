package ws

import (
	"testing"
	"time"

	"github.com/chenjingxiong/weblise/server/protocol"
)

func TestHub_RegisterAgentConnection(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	conn := NewConnection(protocol.NewConnection(protocol.ConnectionTypeAgent, "test-agent-key"))
	hub.RegisterAgentConnection(conn)

	time.Sleep(50 * time.Millisecond)

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

	hub.RegisterAgentConnection(agentConn)
	hub.RegisterClientConnection(clientConn)

	time.Sleep(50 * time.Millisecond)

	session := hub.CreateSession("test-session", clientConn, agentConn)

	if session == nil {
		t.Fatal("Session not created")
	}

	if session.ID != "test-session" {
		t.Errorf("Expected session ID 'test-session', got %s", session.ID)
	}

	if session.AgentConn.ID != agentConn.ID {
		t.Error("Agent connection mismatch")
	}

	if session.ClientConn.ID != clientConn.ID {
		t.Error("Client connection mismatch")
	}
}
