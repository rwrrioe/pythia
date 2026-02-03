package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/auth/authz"
	ocr_grpc_client "github.com/rwrrioe/pythia/backend/internal/clients/ocr/grpc"
	sso_grpc_client "github.com/rwrrioe/pythia/backend/internal/clients/sso/grpc"
	service "github.com/rwrrioe/pythia/backend/internal/services"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	grpcconn "github.com/rwrrioe/pythia/backend/internal/transport/grpc"
	"github.com/rwrrioe/pythia/backend/internal/transport/rest"
	"github.com/rwrrioe/pythia/backend/internal/transport/ws"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type App struct {
	ocrClient    *ocr_grpc_client.Client
	ssoClient    *sso_grpc_client.Client
	wsHandlers   *ws.Handlers
	restHandlers *rest.Handlers
	router       *gin.Engine
}

func New(
	ctx context.Context,
	log *slog.Logger,
	appSecret string,
) (*App, error) {
	const op = "App.New"

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//init storage and repos
	redisClient, err := taskstorage.NewRedisStorage(ctx, "redis:6379", time.Hour)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	pool, err := postgresql.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	deckStorage := postgresql.NewDeckStorage(pool)
	flStorage := postgresql.NewFlashcardStorage(pool)
	ssStorage := postgresql.NewSessionStorage(pool)
	txm := postgresql.NewTxManager(pool)
	//init grpc-clients

	ocrConn, err := grpcconn.New(log, grpcconn.Config{
		Addr:         "ocr:9080",
		Timeout:      time.Minute,
		RetriesCount: 5,
	})
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	ssoConn, err := grpcconn.New(log, grpcconn.Config{
		Addr:         "sso:9081",
		Timeout:      time.Minute,
		RetriesCount: 5,
	})

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	ocrClient := ocr_grpc_client.New(ocrConn, log)
	ssoClient := sso_grpc_client.New(ssoConn, log)

	// init services

	sso := authn.NewSSO(ssoClient, 1)
	ocr := service.NewOCRService(ocrClient)
	learn := service.NewLearnService(4)
	cards := service.NewCardsService()
	transl, err := service.NewTranslateService(ctx, "gemini-2.5-flash-lite")
	stats := service.NewStatsService(ssStorage, deckStorage, flStorage, txm)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	authorizer := authz.NewAuthorizer(redisClient, log)

	session, err := service.NewSessionService(
		ocr,
		transl,
		learn,
		cards,
		redisClient,
		txm,
		ssStorage,
		deckStorage,
		flStorage,
		authorizer,
	)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	// init handlers
	hub := hub.NewWebSocketHub()
	wsHandlers := ws.New(hub)
	ws.RegisterRoutes(router, wsHandlers)
	restHandlers := rest.New(session, stats, sso, hub, redisClient)
	authMiddleware := authn.New(log, appSecret)
	requireAuthMiddleware := authn.NewRequireAuth(log)

	rest.RegisterRoutes(router, restHandlers, authMiddleware(), requireAuthMiddleware())

	return &App{
		ocrClient:    ocrClient,
		ssoClient:    ssoClient,
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
