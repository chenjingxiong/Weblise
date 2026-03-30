package db

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Use in-memory database for testing
	cfg := Config{Path: ":memory:"}
	database, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer database.Close()

	// Check if tables were created
	var tableName string
	err = database.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='devices'").Scan(&tableName)
	if err != nil {
		t.Errorf("Devices table not created: %v", err)
	}
}

func TestCreateDevice(t *testing.T) {
	cfg := Config{Path: ":memory:"}
	database, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer database.Close()

	err = database.CreateDevice("device-1", "test-key", "Test Device", "linux")
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}

	device, err := database.GetDeviceByKey("test-key")
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
	database, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer database.Close()

	device, err := database.GetDeviceByKey("non-existent")
	if err != nil {
		t.Fatalf("GetDeviceByKey failed: %v", err)
	}

	if device != nil {
		t.Error("Expected nil for non-existent device")
	}
}

func TestUpdateDeviceLastSeen(t *testing.T) {
	cfg := Config{Path: ":memory:"}
	database, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer database.Close()

	err = database.CreateDevice("device-1", "test-key", "Test Device", "linux")
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}

	// Wait a bit to ensure timestamp changes
	time.Sleep(10 * time.Millisecond)

	err = database.UpdateDeviceLastSeen("device-1")
	if err != nil {
		t.Fatalf("UpdateDeviceLastSeen failed: %v", err)
	}

	device, err := database.GetDeviceByKey("test-key")
	if err != nil {
		t.Fatalf("GetDeviceByKey failed: %v", err)
	}

	// Check that last seen was updated (should be very recent)
	if time.Since(device.LastSeen) > 5*time.Second {
		t.Error("LastSeen was not updated properly")
	}
}
