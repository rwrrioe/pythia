package rest

import (
	"github.com/gin-gonic/gin"
	rest_handlers "github.com/rwrrioe/pythia/backend/internal/transport/rest/handlers"
)

func RegisterRoutes(r *gin.Engine, handler *rest_handlers.OCRHandler) {
	api := r.Group("/api")
	{
		api.POST("/upload", handler.Upload)
	}
}
