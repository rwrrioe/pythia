package main

import (
	"context"
	"log"
	"time"

	translate "github.com/rwrrioe/pythia/internal/services/translate_service"
)

func main() {
	start := time.Now()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	text := `Die Polizei in Bayern erhielt im letzten Monat die Bewerbung eines 17-jahrigen schillers, komplett mit Anschreiben, Lebenslauf und sehr gutem Zeugnis. Von seinen Unterlagen Uberzeugt, lud das Prasidium in Miinchen den Jungen zu einem Vorstellungsgesprach ein. Ruhig und entspannt konnte er auf alle Fragen des Personalchefs antworten und hinterlie einen sehr guten Eindruck. Nachdem der Abiturient sich verabschiedet hatte, warf der Personalchef einen letzten Blick auf sein Zeugnis — und entdeckte darauf die Jahreszahl 1993! Der Schiller hatte einfach das Abschlusszeugnis seines Vaters verwendet, der den gleichen Namen tragt, und nur an siner Stelle vergessen, das Datum zu korrigieren.`
	res, err := translate.NewTranslateService().FindUnknown(ctx, text)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
	log.Println("taken:", time.Since(start))
}
