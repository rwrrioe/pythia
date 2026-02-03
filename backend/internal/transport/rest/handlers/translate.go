package rest_handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/domain/requests"
	service "github.com/rwrrioe/pythia/backend/internal/services"
	service_errors "github.com/rwrrioe/pythia/backend/internal/services/errors"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type TranslateHandler struct {
	storage *taskstorage.RedisStorage
	ws      *hub.WebSocketHub
	session *service.SessionService
}

func NewTranslateHandler(storage *taskstorage.RedisStorage, ws *hub.WebSocketHub, session *service.SessionService) *TranslateHandler {
	return &TranslateHandler{
		storage: storage,
		ws:      ws,
		session: session,
	}
}

// Translate post /api/session/:sessionId/task/:taskId/translate
func (h *TranslateHandler) Translate(c *gin.Context) {
	sessionId, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session_id is empty",
		})
		return
	}

	taskId := c.Param("taskId")

	var req requests.AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx := c.Request.Context()

	uid, _ := authn.UIDFromContext(ctx)

	bgCtx := context.WithValue(context.Background(), "user_id", uid)
	bgCtx, cancel := context.WithTimeout(bgCtx, 2*time.Minute)
	go func(ctx context.Context) {
		defer cancel()
		h.ws.Notify(sessionId, gin.H{
			"session_id": sessionId,
			"task_id":    taskId,
			"status":     "processing",
			"stage":      "translate",
		})

		words, err := h.session.FindWords(ctx, int64(sessionId), taskId)
		if err != nil {
			if errors.Is(err, service_errors.ErrSessionNotFound) {
				h.ws.Notify(sessionId, gin.H{
					"session_id": sessionId,
					"task_id":    taskId,
					"error":      "session not found",
				})
				return
			}

			if errors.Is(err, service_errors.ErrTaskNotFound) {
				h.ws.Notify(sessionId, gin.H{
					"session_id": sessionId,
					"task_id":    taskId,
					"stage":      "translate",
					"error":      "task not found",
				})
				return
			}

			h.ws.Notify(sessionId, gin.H{
				"session_id": sessionId,
				"task_id":    taskId,
				"stage":      "translate",
				"error":      err.Error(),
			})
			return
		}

		h.ws.Notify(sessionId, gin.H{
			"status":     "done",
			"stage":      "translate",
			"session_id": sessionId,
			"task_id":    taskId,
			"words":      words,
		})
	}(bgCtx)
	c.JSON(http.StatusAccepted, gin.H{
		"session_id": sessionId,
		"task_id":    taskId,
		"stage":      "translate",
	})
}
