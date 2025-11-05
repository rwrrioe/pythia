package rest_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/domain/models"
	translate "github.com/rwrrioe/pythia/backend/internal/services/translate_service"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type TranslateHandler struct {
	ws     *hub.WebSocketHub
	transl *translate.TranslateService
}

func NewTranslateHandler(ws *hub.WebSocketHub, transl *translate.TranslateService) *TranslateHandler {
	return &TranslateHandler{
		ws:     ws,
		transl: transl,
	}
}

// post /api/translate
func (h *TranslateHandler) Translate(c *gin.Context) {
	var req models.AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	go func() {
		h.ws.Notify(req.TaskID, gin.H{
			"status": "processing",
			"stage":  "translate",
		})
		words, err := h.transl.FindUnknownWords(c, req)
		if err != nil {
			h.ws.Notify(req.TaskID, gin.H{"status": "error", "error": err.Error()})
			return
		}
		h.ws.Notify(req.TaskID, gin.H{"status": "done", "words": words})
	}()
	c.JSON(http.StatusAccepted, gin.H{"task_id": string(req.TaskID)})
}
