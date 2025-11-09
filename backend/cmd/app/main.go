package main

import (
	"context"
	"log"
	"os"

	"github.com/rwrrioe/pythia/backend/internal/app"
)

func main() {
	_ = os.Getenv("GEMINI_API_KEY")
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	app, err := app.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	app.MustRun()
}
