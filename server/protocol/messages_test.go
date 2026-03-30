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
