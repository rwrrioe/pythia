package rest_handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/rwrrioe/pythia/backend/internal/services"
)

type StatsHandler struct {
	stats *service.StatsService
}

func NewStatsHandler(stats *service.StatsService) *StatsHandler {
	return &StatsHandler{
		stats: stats,
	}
}

// GET /api/dashboard
func (h *StatsHandler) Dashboard(c *gin.Context) {
	ctx := c.Request.Context()

	dashboard, err := h.stats.Dashboard(ctx)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "user is unauthorized",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dashboard": dashboard,
	})
}
