package ws_hub

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketHub struct {
	mu       sync.Mutex
	sessions map[uuid.UUID]map[*websocket.Conn]struct{}
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{sessions: make(map[uuid.UUID]map[*websocket.Conn]struct{})}
}

func (h *WebSocketHub) Add(sessionId uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.sessions[sessionId] == nil {
		h.sessions[sessionId] = make(map[*websocket.Conn]struct{})
	}

	h.sessions[sessionId][conn] = struct{}{}
}

func (h *WebSocketHub) Notify(sessionId uuid.UUID, payload interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn, _ := range h.sessions[sessionId] {
		_ = conn.WriteJSON(payload)
	}
}

func (h *WebSocketHub) Remove(sessionId uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conns, ok := h.sessions[sessionId]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.sessions, sessionId)
		}
	}
}
