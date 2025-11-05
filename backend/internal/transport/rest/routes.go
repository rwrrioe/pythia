package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/services"
	rest_handlers "github.com/rwrrioe/pythia/backend/internal/transport/rest/handlers"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type Handlers struct {
	OCR      *rest_handlers.OCRHandler
	Transl   *rest_handlers.TranslateHandler
	Services *services.Services
}

func New(services *services.Services, hub *hub.WebSocketHub) *Handlers {
	ocr := rest_handlers.NewOCRHandler(services.OCR, hub)
	transl := rest_handlers.NewTranslateHandler(hub, services.Translate)

	return &Handlers{
		Services: services,
		OCR:      ocr,
		Transl:   transl,
	}
}

func RegisterRoutes(r *gin.Engine, handlers *Handlers) {
	api := r.Group("/api")
	{
		api.POST("/upload", handlers.OCR.Upload)
		api.POST("/translate", handlers.Transl.Translate)
	}
}
