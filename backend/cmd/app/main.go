package main

import (
	"github.com/rwrrioe/pythia/backend/internal/app"
	"github.com/rwrrioe/pythia/backend/internal/services/ocr_service/ocr"
	rest_handlers "github.com/rwrrioe/pythia/backend/internal/transport/rest/handlers"
	ws_handlers "github.com/rwrrioe/pythia/backend/internal/transport/ws/handlers"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

func main() {
	ocrProcessor, err := ocr.NewOCRProcessor()
	if err != nil {
		panic(err)
	}
	defer ocrProcessor.Close()

	wsHub := hub.NewWebSocketHub()

	wsHandler := ws_handlers.NewWebSocketHandler(ocrProcessor, wsHub)
	restHandler := rest_handlers.NewOCRHandler(ocrProcessor, wsHub)

	app := app.New(
		ocrProcessor,
		wsHandler,
		restHandler,
	)
	app.MustRun()
}
