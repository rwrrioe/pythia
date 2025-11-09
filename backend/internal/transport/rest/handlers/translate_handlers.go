package rest_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/domain/models"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/services/task_storage"
	translate "github.com/rwrrioe/pythia/backend/internal/services/translate_service"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type TranslateHandler struct {
	storage *taskstorage.RedisTaskStorage
	ws      *hub.WebSocketHub
	transl  *translate.TranslateService
}

func NewTranslateHandler(ws *hub.WebSocketHub, transl *translate.TranslateService, storage *taskstorage.RedisTaskStorage) *TranslateHandler {
	return &TranslateHandler{
		ws:      ws,
		transl:  transl,
		storage: storage,
	}
}

func (h *TranslateHandler) isTranslateErr(ok bool, err error, req models.AnalyzeRequest, stage string) bool {
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

// post /api/translate
func (h *TranslateHandler) Translate(c *gin.Context) {
	var req models.AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	go func() {
		h.ws.Notify(req.TaskID, gin.H{
			"status": "processing",
			"stage":  "translate",
		})
		task, ok, err := h.storage.Get(c, req.TaskID)
		if h.isTranslateErr(ok, err, req, "translate") {
			return
		}

		words, err := h.transl.FindUnknownWords(c, task, req)
		if h.isTranslateErr(true, err, req, "translate") {
			return
		}

		ok, err = h.storage.UpdateTask(c, req.TaskID, func(task *taskstorage.TaskDTO) {
			task.Words = words
		})
		if h.isTranslateErr(ok, err, req, "translate") {
			return
		}

		h.ws.Notify(req.TaskID, gin.H{"status": "done", "words": words, "stage": "translate"})
	}()
	c.JSON(http.StatusAccepted, gin.H{"task_id": string(req.TaskID), "stage": "translate"})
}

// post /api/translate/examples
func (h *TranslateHandler) WriteExamples(c *gin.Context) {
	var req models.AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	go func() {
		h.ws.Notify(req.TaskID, gin.H{
			"status": "processing",
			"stage":  "writing examples",
		})
		task, ok, err := h.storage.Get(c, req.TaskID)
		if h.isTranslateErr(ok, err, req, "writing examples") {
			return
		}

		examples, err := h.transl.WriteExamples(c, task, req)
		if h.isTranslateErr(true, err, req, "writing examples") {
			return
		}

		ok, err = h.storage.UpdateTask(c, req.TaskID, func(task *taskstorage.TaskDTO) {
			task.Examples = examples
		})
		if h.isTranslateErr(ok, err, req, "writing examples") {
			return
		}

		h.ws.Notify(req.TaskID, gin.H{"status": "done", "words": examples, "stage": "writing examples"})
	}()
	c.JSON(http.StatusAccepted, gin.H{"task_id": string(req.TaskID), "stage": "writing examples"})
}
