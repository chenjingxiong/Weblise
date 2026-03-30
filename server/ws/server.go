package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 32,
	WriteBufferSize: 1024 * 32,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Server manages the WebSocket server
type Server struct {
	addr string
	hub  *Hub
}

// NewServer creates a new WebSocket server
func NewServer(addr string, hub *Hub) *Server {
	return &Server{
		addr: addr,
		hub:  hub,
	}
}

// Start starts the WebSocket server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/agent", s.handleAgentConnection)
	mux.HandleFunc("/client", s.handleClientConnection)

	log.Printf("[WS Server] Starting on %s", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) handleAgentConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS Server] Agent upgrade failed: %v", err)
		return
	}
	log.Printf("[WS Server] New agent connection from %s", r.RemoteAddr)
	go s.handleAgent(conn)
}

func (s *Server) handleClientConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS Server] Client upgrade failed: %v", err)
		return
	}
	log.Printf("[WS Server] New client connection from %s", r.RemoteAddr)
	go s.handleClient(conn)
}
