package ws

import (
	"github.com/gin-gonic/gin"
	ws_handlers "github.com/rwrrioe/pythia/backend/internal/transport/ws/handlers"
)

func RegisterRoutes(r *gin.Engine, handler *ws_handlers.WebSocketOCRHandler) {
	ws := r.Group("/ws")
	{
		ws.GET("/ocr", handler.WebSocket)
	}
}
