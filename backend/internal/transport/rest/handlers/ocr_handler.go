package rest_handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/services/ocr_service/ocr"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type OCRHandler struct {
	ocr *ocr.OCRProcesser
	ws  *hub.WebSocketHub
}

func NewOCRHandler(ocr *ocr.OCRProcesser, ws *hub.WebSocketHub) *OCRHandler {
	return &OCRHandler{ocr: ocr, ws: ws}
}

func (h *OCRHandler) Upload(c *gin.Context) {
	taskID := c.PostForm("task_id")
	if taskID == "" {
		taskID = uuid.NewString()
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while fileheader open": err.Error()})
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while io read all": err.Error()})
		return
	}

	go func() {
		h.ws.Notify(taskID, gin.H{"status": "processing"})

		texts, err := h.ocr.ProcessImage(c, data, "de")
		if err != nil {
			h.ws.Notify(taskID, gin.H{"status": "error", "error": err.Error()})
			return
		}

		h.ws.Notify(taskID, gin.H{"status": "done", "rec_texts": texts})
	}()
	c.JSON(http.StatusAccepted, gin.H{"task_id": taskID})
}
