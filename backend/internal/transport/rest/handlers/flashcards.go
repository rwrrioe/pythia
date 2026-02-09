package rest_handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/rwrrioe/pythia/backend/internal/services"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type FlashCardsHandler struct {
	storage *taskstorage.RedisStorage
	session *service.SessionService
	ws      *hub.WebSocketHub
}

func NewFlashCardsHandler(storage *taskstorage.RedisStorage, ws *hub.WebSocketHub, session *service.SessionService) *FlashCardsHandler {
	return &FlashCardsHandler{
		storage: storage,
		session: session,
		ws:      ws,
	}
}

// FlashCards post /api/session/:sessionId/flashcards
func (h *FlashCardsHandler) FlashCards(c *gin.Context) {

	sessionId, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session_id is empty",
		})
		return
	}

	ctx := c.Request.Context()
	flCards, err := h.session.GetFlashcards(ctx, int64(sessionId))
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "session not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionId,
		"stage":      "flashcards",
		"flashcards": flCards,
	})
}
