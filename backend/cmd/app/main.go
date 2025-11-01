package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/services/ocr_service/ocr"
	"github.com/rwrrioe/pythia/backend/internal/transport/ws"
	ws_handlers "github.com/rwrrioe/pythia/backend/internal/transport/ws/handlers"
)

func main() {
	ocrClient, err := ocr.NewOCRProcessor("localhost:50051")
	if err != nil {
		log.Print(err.Error())
	}

	r := gin.Default()
	hub := ws.NewWebSocketHub()
	wsHandler := ws_handlers.NewWebSocketHandler(ocrClient, hub)
	r.GET("/ws/ocr", wsHandler.WebSocket)
	r.POST("/api/ocr", wsHandler.Upload)
}
