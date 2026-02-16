package rest_handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/domain/requests"
	service "github.com/rwrrioe/pythia/backend/internal/services"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type SessionHandler struct {
	storage *taskstorage.RedisStorage
	session *service.SessionService
	ws      *hub.WebSocketHub
}

func NewSessionHandler(storage *taskstorage.RedisStorage, ws *hub.WebSocketHub, session *service.SessionService) *SessionHandler {

	return &SessionHandler{
		storage: storage,
		ws:      ws,
		session: session,
	}
}

// /api/session/new
func (h *SessionHandler) NewSession(c *gin.Context) {
	var req requests.CreateSession

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Printf("new session req lang %d", req.LangId)

	ctx := c.Request.Context()

	id, err := h.session.StartSession(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": id,
	})
}

// /api/session/:sessionId/end
func (h *SessionHandler) EndSession(c *gin.Context) {
	sessionId, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid sessionId",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	err = h.session.EndSession(ctx, sessionId)

	if err != nil && errors.Is(err, service.ErrUnauthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user is unauthorized",
			"details": err.Error(),
		})
		return
	}

	if err != nil && errors.Is(err, service.ErrForbidden) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "access forbidden",
			"details": err.Error(),
		})
		return
	}

	if err != nil && errors.Is(err, service.ErrSessionAlreadyFinished) {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "session already finished",
			"details": err.Error(),
		})
		return
	}
	if err != nil && errors.Is(err, service.ErrSessionNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "session not found",
			"details": err.Error(),
		})
		return
	}

	if err != nil && errors.Is(err, service.ErrNoWords) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "no words",
			"details": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal error",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionId,
		"status":     "finished",
	})
}

// /api/session/:sessionId/summary
func (h *SessionHandler) SessionSummary(c *gin.Context) {
	var req requests.SummarizeSession

	sessionId, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid sessionId",
			"details": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	//summarize session
	words, err := h.session.SummarizeSession(ctx, sessionId, req.Accuracy)
	if err != nil && errors.Is(err, service.ErrUnauthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user is unauthorized",
			"details": err.Error(),
		})
		return
	}

	if err != nil && errors.Is(err, service.ErrForbidden) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "access forbidden",
			"details": err.Error(),
		})
		return
	}

	if err != nil && errors.Is(err, service.ErrSessionNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "session not found",
			"details": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"words":       words,
		"words_count": len(words),
		"accuracy":    req.Accuracy,
	})

}
