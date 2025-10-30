package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rwrrioe/pythia/backend/internal/domain/models"
	translate "github.com/rwrrioe/pythia/backend/internal/services/translate_service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}
	start := time.Now()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	text := `Die Polizei in Bayern erhielt im letzten Monat die Bewerbung eines 17-jahrigen schillers, komplett mit Anschreiben, Lebenslauf und sehr gutem Zeugnis. Von seinen Unterlagen Uberzeugt, lud das Prasidium in Miinchen den Jungen zu einem Vorstellungsgesprach ein. Ruhig und entspannt konnte er auf alle Fragen des Personalchefs antworten und hinterlie einen sehr guten Eindruck. Nachdem der Abiturient sich verabschiedet hatte, warf der Personalchef einen letzten Blick auf sein Zeugnis — und entdeckte darauf die Jahreszahl 1993! Der Schiller hatte einfach das Abschlusszeugnis seines Vaters verwendet, der den gleichen Namen tragt, und nur an siner Stelle vergessen, das Datum zu korrigieren.`

	req := models.AnalyzeRequest{
		Text:     []byte(text),
		Level:    "B1",
		Durating: "6 months",
		Book:     "Schritte neu 5 International",
		Lang:     "DE",
	}

	service, err := translate.NewTranslateService(ctx, "gemini-2.5-flash")
	if err != nil {
		log.Fatal(err)
	}
	res, err := service.FindUnknownWords(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
	log.Println("taken:", time.Since(start))
}
