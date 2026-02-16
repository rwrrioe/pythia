package rest_handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	service "github.com/rwrrioe/pythia/backend/internal/services"
	storage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type OCRHandler struct {
	storage *storage.RedisStorage
	session *service.SessionService
	ws      *hub.WebSocketHub
}

func NewOCRHandler(storage *storage.RedisStorage, ws *hub.WebSocketHub, session *service.SessionService) *OCRHandler {
	return &OCRHandler{
		session: session,
		ws:      ws,
		storage: storage,
	}
}

// /!!!!
func (h *OCRHandler) respondOCRErr(c *gin.Context, err error, code int, message string) {
	c.JSON(code, gin.H{message: err.Error()})
	c.Abort()
}

// /api/session/:sessionId/upload
func (h *OCRHandler) Upload(c *gin.Context) {
	sessionId, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		h.respondOCRErr(c, err, http.StatusBadRequest, "invalid sessionId")
		return
	}

	taskID := c.PostForm("task_id")
	if taskID == "" {
		taskID = uuid.NewString()
	}

	lang := c.PostForm("lang")
	if lang == "" {
		h.respondOCRErr(c, fmt.Errorf("no language"), http.StatusBadRequest, "no language")
		return
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
		return
	}

	ctx := c.Request.Context()

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user is unauthorized",
			"details": "",
		})
		return
	}

	//creating background ctx !!TODO -> implement on gin
	bgCtx := context.WithValue(context.Background(), "user_id", uid)
	bgCtx, cancel := context.WithTimeout(bgCtx, 2*time.Minute)
	go func(ctx context.Context) {
		defer cancel()
		h.ws.Notify(sessionId, gin.H{
			"task_id":    taskID,
			"session_id": sessionId,
			"status":     "processing",
			"stage":      "ocr",
		})
		err := h.session.RecognizeText(ctx, sessionId, taskID, data, lang)
		if err != nil {
			h.ws.Notify(sessionId, gin.H{
				"task_id":    taskID,
				"session_id": sessionId,
				"status":     "error",
				"error":      err.Error(),
				"stage":      "ocr"})
			return
		}

		h.ws.Notify(sessionId, gin.H{
			"task_id":    taskID,
			"session_id": sessionId,
			"status":     "done",
			"stage":      "ocr"})
	}(bgCtx)
	c.JSON(http.StatusAccepted, gin.H{
		"task_id":    taskID,
		"session_id": sessionId,
		"stage":      "ocr"})
}
