package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/chenjingxiong/weblise/server/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (s *Server) handleAgent(wsConn *websocket.Conn) {
	defer wsConn.Close()

	var msg protocol.Message
	if err := wsConn.ReadJSON(&msg); err != nil {
		log.Printf("[Agent] Read register failed: %v", err)
		return
	}

	if msg.Type != protocol.MessageTypeRegister {
		log.Printf("[Agent] First message must be register")
		return
	}

	var regMsg protocol.RegisterMessage
	if err := msg.ParseData(&regMsg); err != nil {
		log.Printf("[Agent] Parse register failed: %v", err)
		return
	}

	log.Printf("[Agent] Register: key=%s, name=%s", regMsg.Key, regMsg.Name)

	baseConn := protocol.NewConnection(protocol.ConnectionTypeAgent, regMsg.Key)
	conn := NewConnection(baseConn)

	s.hub.RegisterAgentConnection(conn.Connection)

	done := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-conn.SendChan:
				data, _ := json.Marshal(msg)
				if err := wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
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

	for {
		_, data, err := wsConn.ReadMessage()
		if err != nil {
			break
		}

		conn.LastSeen = time.Now()

		var msg protocol.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case protocol.MessageTypeHeartbeat:
			conn.SendSafe(&protocol.Message{Type: protocol.MessageTypePong})
		default:
			s.hub.messageAgent <- &AgentMessage{ConnID: conn.ID, Msg: &msg}
		}
	}

	close(done)
	s.hub.Unregister(conn.Connection)
}
