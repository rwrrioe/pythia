package rest_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/domain/models"
	learn "github.com/rwrrioe/pythia/backend/internal/services/learn_service"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type LearnHandler struct {
	storage *taskstorage.RedisTaskStorage
	ws      *hub.WebSocketHub
	learn   *learn.LearnService
}

func NewLearnHandler(ws *hub.WebSocketHub, learn *learn.LearnService, storage *taskstorage.RedisTaskStorage) *LearnHandler {
	return &LearnHandler{
		ws:      ws,
		learn:   learn,
		storage: storage,
	}
}

func (h *LearnHandler) isErr(ok bool, err error, req models.AnalyzeRequest, stage string) bool {
	if ok != true {
		h.ws.Notify(req.TaskID, gin.H{"status": "not found", "stage": stage})
		return true
	}
	if err != nil {
		h.ws.Notify(req.TaskID, gin.H{"status": "error", "error": err, "stage": stage})
		return true
	}
	return false
}

// post /api/flashcards/tests/quiz
func (h *LearnHandler) Quiz(c *gin.Context) {
	var req models.AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	///!!!!
	go func() {
		h.ws.Notify(req.TaskID, gin.H{
			"status": "processing",
			"stage":  "quiz",
		})
		task, ok, err := h.storage.Get(c, req.TaskID)
		if h.isErr(ok, err, req, "quiz") {
			return
		}

		quiz := h.learn.QuizTest(c, &task.Words)
		h.ws.Notify(req.TaskID, gin.H{"status": "done", "quiz": quiz, "stage": "quiz"})
	}()

	c.JSON(http.StatusAccepted, gin.H{"task_id": req.TaskID, "stage": "quiz"})
}
