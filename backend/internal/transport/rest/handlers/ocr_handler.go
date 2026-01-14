package rest_handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/services/ocr_service/ocr"
	storage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type OCRHandler struct {
	storage *storage.RedisTaskStorage
	ocr     *ocr.OCRProcesser
	ws      *hub.WebSocketHub
}

func NewOCRHandler(ocr *ocr.OCRProcesser, ws *hub.WebSocketHub, storage *storage.RedisTaskStorage) *OCRHandler {
	return &OCRHandler{
		ocr:     ocr,
		ws:      ws,
		storage: storage,
	}
}

// /!!!!
func (h *OCRHandler) respondOCRErr(c *gin.Context, err error, code int, message string) {
	c.JSON(code, gin.H{message: err.Error()})
	c.Abort()
}

func (h *OCRHandler) Upload(c *gin.Context) {
	taskID := c.PostForm("task_id")
	if taskID == "" {
		taskID = uuid.NewString()
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.respondOCRErr(c, err, http.StatusBadRequest, "no file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.respondOCRErr(c, err, http.StatusInternalServerError, "can't open file")
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			h.respondOCRErr(c, err, http.StatusInternalServerError, "can't close file")
			return
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		h.respondOCRErr(c, err, http.StatusInternalServerError, "error while reading file")
	}

	go func() {
		h.ws.Notify(taskID, gin.H{
			"status": "processing",
			"stage":  "ocr",
		})
		texts, err := h.ocr.ProcessImage(c, data, "de")
		if err != nil {
			h.ws.Notify(taskID, gin.H{"status": "error", "error": err.Error(), "stage": "ocr"})
			return
		}

		err = h.storage.Save(c, taskID, &storage.TaskDTO{OCRText: texts})
		if err != nil {
			h.ws.Notify(taskID, gin.H{"status": "error", "error": err.Error(), "stage": "ocr"})
			return
		}

		h.ws.Notify(taskID, gin.H{"status": "done", "stage": "ocr"})
	}()
	c.JSON(http.StatusAccepted, gin.H{"task_id": taskID, "stage": "ocr"})
}
