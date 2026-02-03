package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/services"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	rest_handlers "github.com/rwrrioe/pythia/backend/internal/transport/rest/handlers"
	hub "github.com/rwrrioe/pythia/backend/internal/transport/ws/ws_hub"
)

type Handlers struct {
	authHandler       *rest_handlers.AuthHandler
	ocrHandler        *rest_handlers.OCRHandler
	translateHandler  *rest_handlers.TranslateHandler
	learnHandler      *rest_handlers.LearnHandler
	sessionHandler    *rest_handlers.SessionHandler
	flashcardsHandler *rest_handlers.FlashCardsHandler
	statsHandler      *rest_handlers.StatsHandler
}

func New(session *service.SessionService, stats *service.StatsService, sso authn.SSOService, ws *hub.WebSocketHub, storage *taskstorage.RedisStorage) *Handlers {

	ocr := rest_handlers.NewOCRHandler(storage, ws, session)
	transl := rest_handlers.NewTranslateHandler(storage, ws, session)
	flCards := rest_handlers.NewFlashCardsHandler(storage, ws, session)
	learn := rest_handlers.NewLearnHandler(storage, ws, session)
	ss := rest_handlers.NewSessionHandler(storage, ws, session)
	authH := rest_handlers.NewAuthHandler(sso)
	statsH := rest_handlers.NewStatsHandler(stats)

	return &Handlers{
		ocrHandler:        ocr,
		translateHandler:  transl,
		learnHandler:      learn,
		sessionHandler:    ss,
		authHandler:       authH,
		flashcardsHandler: flCards,
		statsHandler:      statsH,
	}
}

func RegisterRoutes(r *gin.Engine, handlers *Handlers, auth gin.HandlerFunc, requireAuth gin.HandlerFunc) {
	api := r.Group("/api")
	api.Use(auth)

	protected := api.Group("")
	protected.Use(requireAuth)
	protected.GET("/dashboard", handlers.statsHandler.Dashboard)

	session := api.Group("/session")
	session.POST("/new", handlers.sessionHandler.NewSession)

	sessionProtected := session.Group("")
	sessionProtected.Use(requireAuth)
	{
		sessionProtected.POST("/:sessionId/upload", handlers.ocrHandler.Upload)
		sessionProtected.POST("/:sessionId/task/:taskId/translate", handlers.translateHandler.Translate)
		sessionProtected.PATCH("/:sessionId/end", handlers.sessionHandler.EndSession)
		sessionProtected.GET("/:sessionId/learn/flashcards", handlers.flashcardsHandler.FlashCards)
		sessionProtected.GET("/:sessionId/learn/quiz", handlers.learnHandler.Quiz)
		sessionProtected.POST("/:sessionId/summary", handlers.sessionHandler.SessionSummary)
	}

	public := r.Group("/api/auth")
	{
		public.POST("/login", handlers.authHandler.Login)
		public.POST("/register", handlers.authHandler.Register)
	}
}
