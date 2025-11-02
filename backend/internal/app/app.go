package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/domain"
	"github.com/rwrrioe/pythia/backend/internal/transport/rest"
	rest_handlers "github.com/rwrrioe/pythia/backend/internal/transport/rest/handlers"
	"github.com/rwrrioe/pythia/backend/internal/transport/ws"
	ws_handlers "github.com/rwrrioe/pythia/backend/internal/transport/ws/handlers"
)

type App struct {
	ocr domain.ImageProcesser
	// transl     domain.TranslateProvider
	// flashCards domain.CardsBuilder
	router *gin.Engine
}

func New(ocr domain.ImageProcesser, wsHandler *ws_handlers.WebSocketOCRHandler, restHandler *rest_handlers.OCRHandler) *App {
	router := gin.Default()
	rest.RegisterRoutes(router, restHandler)
	ws.RegisterRoutes(router, wsHandler)
	return &App{
		ocr: ocr,
		// transl:     transl,
		// flashCards: flashCards,
		router: router,
	}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}

func (a *App) MustRun(addr string) {
	err := a.Run(addr)
	if err != nil {
		panic(err)
	}
}
