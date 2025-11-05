package ws_hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketHub struct {
	mu      sync.Mutex
	clients map[string]*websocket.Conn
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{clients: make(map[string]*websocket.Conn)}
}

func (h *WebSocketHub) AddClient(id string, conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[id] = conn
	h.mu.Unlock()
}

func (h *WebSocketHub) Notify(id string, payload interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, ok := h.clients[id]; ok {
		if err := conn.WriteJSON(payload); err != nil {
			log.Printf("Failed to send WS message to %s: %v", id, err)

			conn.Close()
			delete(h.clients, id)
		}
	}
}

func (h *WebSocketHub) RemoveClient(id string) {
	h.mu.Lock()
	delete(h.clients, id)
	h.mu.Unlock()
}
