package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/services"
	ws_handlers "github.com/rwrrioe/pythia/backend/internal/transport/ws/handlers"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type Handlers struct {
	Services *services.Services
	OCR      *ws_handlers.OCRHandler
	Transl   *ws_handlers.TranslateHandler
}

func New(services *services.Services, ws *hub.WebSocketHub) *Handlers {
	ocr := ws_handlers.NewOCRHandler(services.OCR, ws)
	transl := ws_handlers.NewTranslateHandler(services.Translate, ws)

	return &Handlers{
		Services: services,
		OCR:      ocr,
		Transl:   transl,
	}
}

func RegisterRoutes(r *gin.Engine, handlers *Handlers) {
	ws := r.Group("/ws")
	{
		ws.GET("/ocr", handlers.OCR.WebSocket)
		ws.GET("/translate", handlers.Transl.WebSocket)
	}
}
