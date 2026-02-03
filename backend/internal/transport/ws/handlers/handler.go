package ws_handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	ws *hub.WebSocketHub
}

func NewHandler(ws *hub.WebSocketHub) *Handler {
	return &Handler{
		ws: ws,
	}
}

func (h *Handler) WebSocket(c *gin.Context) {
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid sessionId",
			"details": err.Error(),
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("failed to upgrade")
		return
	}

	h.ws.Add(sessionId, conn)
}
