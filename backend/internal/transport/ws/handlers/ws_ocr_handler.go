package ws_handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rwrrioe/pythia/backend/internal/services/ocr_service/ocr"
	"github.com/rwrrioe/pythia/backend/internal/transport/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketOCRHandler struct {
	ocr *ocr.OCRProcessor
	ws  *ws.WebSocketHub
}

// get /ws/ocr
func NewWebSocketHandler(ocr *ocr.OCRProcessor, ws *ws.WebSocketHub) *WebSocketOCRHandler {
	return &WebSocketOCRHandler{ocr: ocr, ws: ws}
}

func (h *WebSocketOCRHandler) WebSocket(c *gin.Context) {
	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing task_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	h.ws.AddClient(taskID, conn)
}

// post /api/ocr
func (h *WebSocketOCRHandler) Upload(c *gin.Context) {
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

		texts, err := h.ocr.RecognizeText(c, data, "de")
		if err != nil {
			h.ws.Notify(taskID, gin.H{"status": "error", "error": err.Error()})
			return
		}

		h.ws.Notify(taskID, gin.H{"status": "done", "rec_texts": texts})
	}()
	c.JSON(http.StatusAccepted, gin.H{"task_id": taskID})
}
