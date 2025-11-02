package ws_handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rwrrioe/pythia/backend/internal/services/ocr_service/ocr"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketOCRHandler struct {
	ocr *ocr.OCRProcesser
	ws  *hub.WebSocketHub
}

func NewWebSocketHandler(ocr *ocr.OCRProcesser, ws *hub.WebSocketHub) *WebSocketOCRHandler {
	return &WebSocketOCRHandler{
		ocr: ocr,
		ws:  ws,
	}
}

func (h WebSocketOCRHandler) WebSocket(c *gin.Context) {
	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing task_id",
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("failed to upgrade")
		return
	}

	h.ws.AddClient(taskID, conn)
}
