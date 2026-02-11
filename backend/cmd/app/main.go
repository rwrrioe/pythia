package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/rwrrioe/pythia/backend/internal/app"
	config "github.com/rwrrioe/pythia/backend/internal/config/grpconn"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	_ = os.Getenv("GEMINI_API_KEY")
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	env := os.Getenv("LOGGER_ENV")
	appSecret := os.Getenv("APP_SECRET")
	log := setupLogger(env)
	log.Info("starting app", slog.Any("env", env))

	ocrCfg, err := config.FetchConfig(config.ConfigAttr{
		CfgType: config.OCR,
	})
	if err != nil {
		log.Error("failed to fetch ocr config")
		panic("failed to fetch ocr config")
	}

	ssoCfg, err := config.FetchConfig(config.ConfigAttr{
		CfgType: config.SSO,
	})
	if err != nil {
		log.Error("failed to fetch ocr config")
		panic("failed to fetch ocr config")
	}

	app, err := app.New(ctx, log, appSecret, ssoCfg, ocrCfg)
	if err != nil {
		panic(err)
	}
	app.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
