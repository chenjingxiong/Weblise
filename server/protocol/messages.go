package protocol

import "encoding/json"

// MessageType represents the type of message
type MessageType string

const (
	// Agent messages
	MessageTypeRegister  MessageType = "register"
	MessageTypeHeartbeat MessageType = "ping"
	MessageTypeError     MessageType = "error"
	MessageTypeFrame     MessageType = "frame"

	// Client messages
	MessageTypeConnect    MessageType = "connect"
	MessageTypeDisconnect MessageType = "disconnect"
	MessageTypeInput      MessageType = "input"

	// Bidirectional
	MessageTypePong MessageType = "pong"
)

// Message is the base message structure
type Message struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// RegisterMessage is sent by agent to register with server
type RegisterMessage struct {
	Key    string `json:"key"`
	Name   string `json:"name,omitempty"`
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
	Action string    `json:"action"` // mousemove, mousedown, mouseup, keydown, keyup
	Data   InputData `json:"data"`
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
