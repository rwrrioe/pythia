package rest_handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/rwrrioe/pythia/backend/internal/services"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type LearnHandler struct {
	storage *taskstorage.RedisStorage
	ws      *hub.WebSocketHub
	session *service.SessionService
}

func NewLearnHandler(storage *taskstorage.RedisStorage, ws *hub.WebSocketHub, session *service.SessionService) *LearnHandler {
	return &LearnHandler{
		ws:      ws,
		session: session,
		storage: storage,
	}
}

// post /api/session/:sessionId/quiz
func (h *LearnHandler) Quiz(c *gin.Context) {

	sessionId, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session_id is empty",
		})
		return
	}

	ctx := c.Request.Context()
	quiz, err := h.session.Quiz(ctx, int64(sessionId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"session_id": sessionId,
		"stage":      "quiz",
		"quiz":       quiz,
	})
}
