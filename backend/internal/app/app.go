package app

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/services"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/services/task_storage"
	"github.com/rwrrioe/pythia/backend/internal/transport/rest"
	"github.com/rwrrioe/pythia/backend/internal/transport/ws"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type App struct {
	services     *services.Services
	wsHandlers   *ws.Handlers
	restHandlers *rest.Handlers
	router       *gin.Engine
}

func New(ctx context.Context) (*App, error) {
	router := gin.Default()
	services, err := services.New(ctx, "gemini-2.5-flash-lite", "ocr:9080")
	if err != nil {
		return nil, err
	}
	redisClient := taskstorage.NewRedisTaskStorage("redis:6379", time.Hour)
	hub := hub.NewWebSocketHub()
	wsHandlers := ws.New(services, hub)
	ws.RegisterRoutes(router, wsHandlers)
	restHandlers := rest.New(services, hub, redisClient)
	rest.RegisterRoutes(router, restHandlers)

	return &App{
		services:     services,
		wsHandlers:   wsHandlers,
		restHandlers: restHandlers,
		router:       router,
	}, nil
}

func (a *App) Run() error {
	return a.router.Run()
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
