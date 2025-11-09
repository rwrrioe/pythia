package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/services"
	ws_handlers "github.com/rwrrioe/pythia/backend/internal/transport/ws/handlers"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type Handlers struct {
	Services *services.Services
	Ws       *ws_handlers.Handler
}

func New(services *services.Services, ws *hub.WebSocketHub) *Handlers {
	h := ws_handlers.NewHandler(services.OCR, ws)

	return &Handlers{
		Services: services,
		Ws:       h,
	}
}

func RegisterRoutes(r *gin.Engine, handlers *Handlers) {
	r.GET("/ws", handlers.Ws.WebSocket)
}
