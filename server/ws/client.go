package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/chenjingxiong/weblise/server/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (s *Server) handleClient(wsConn *websocket.Conn) {
	defer wsConn.Close()

	baseConn := protocol.NewConnection(protocol.ConnectionTypeClient, "")
	conn := NewConnection(baseConn)

	s.hub.RegisterClientConnection(conn.Connection)

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
		case protocol.MessageTypeConnect:
			var connMsg protocol.ConnectMessage
			if err := msg.ParseData(&connMsg); err != nil {
				continue
			}

			agentConn := s.hub.GetAgentByKey(connMsg.AgentKey)
			if agentConn == nil {
				conn.SendSafe(&protocol.Message{Type: protocol.MessageTypeError})
				continue
			}

			sessionID := uuid.New().String()
			s.hub.CreateSession(sessionID, conn.Connection, agentConn)
			conn.SendSafe(&protocol.Message{Type: protocol.MessageTypeConnect})

		case protocol.MessageTypeHeartbeat:
			conn.SendSafe(&protocol.Message{Type: protocol.MessageTypePong})

		default:
			s.hub.messageClient <- &ClientMessage{ConnID: conn.ID, Msg: &msg}
		}
	}

	close(done)
	s.hub.Unregister(conn.Connection)
}
