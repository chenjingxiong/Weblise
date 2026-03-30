package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

const serverAddr = "http://localhost:8080"

func TestMain(m *testing.M) {
	if err := startServer(); err != nil {
		os.Exit(1)
	}

	if !waitForServer(serverAddr + "/health") {
		os.Exit(1)
	}

	code := m.Run()
	stopServer()
	os.Exit(code)
}

func startServer() error {
	cmd := exec.Command("go", "run", "../server/main.go")
	cmd.Dir = ".."
	return cmd.Start()
}

func stopServer() {
	// Kill server process
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

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result["status"])
	}
}
