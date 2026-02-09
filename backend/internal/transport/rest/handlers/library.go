package rest_handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	service "github.com/rwrrioe/pythia/backend/internal/services"
)

type LibraryHandler struct {
	library    *service.LibraryService
	flashcards *service.FlashCardsService
	log        *slog.Logger
}

func NewLibraryHandler(library *service.LibraryService, flashcards *service.FlashCardsService, log *slog.Logger) *LibraryHandler {
	return &LibraryHandler{
		library:    library,
		flashcards: flashcards,
		log:        log,
	}
}

// GET /api/library/sessions
func (h *LibraryHandler) ListSession(c *gin.Context) {

	ctx := c.Request.Context()
	ss, err := h.library.Library(ctx)

	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": ss,
	})
}

// GET /api/library/session/:sessionId
func (h *LibraryHandler) GetSession(c *gin.Context) {
	sessionId, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	ctx := c.Request.Context()
	session, err := h.library.GetSession(ctx, int64(sessionId))

	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})

			if errors.Is(err, service.ErrSessionNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "session not found",
					"details": err.Error(),
				})
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}
	flashcards, err := h.flashcards.GetBySession(ctx, int64(sessionId))

	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})

			if errors.Is(err, service.ErrDeckNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "deck not found",
					"details": err.Error(),
				})
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

	}

	h.log.Info("flashcards are ", flashcards)

	var dtos []entities.FlashCardDTO
	for _, fl := range flashcards {
		dtos = append(dtos, entities.FlashCardDTO{
			Word:        fl.Word,
			Translation: fl.Transl,
			Lang:        service.LangsMap[fl.Lang],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"flashcards": dtos,
		"session":    session,
	})
}
