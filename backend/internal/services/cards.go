package service

import (
	"context"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
)

type FlashCardProvider interface {
	List(ctx context.Context) ([]entities.FlashCard, error)
	ListByDeck(ctx context.Context, deckId int) ([]entities.FlashCard, error)
	SaveFlashcard(ctx context.Context, flCard entities.FlashCard) error
}

type DeckProvider interface {
	ListBySession(ctx context.Context, sessionId int) (*entities.Deck, error)
	AttachFlashcard(ctx context.Context, deckId int, flId int) error
	SaveDeck(ctx context.Context, deck entities.Deck) error
}

type FlashCardsService struct{}

func NewCardsService() *FlashCardsService {
	return &FlashCardsService{}
}

func (c *FlashCardsService) BuildCards(ctx context.Context, words []entities.Word) []entities.FlashCardDTO {
	dto := make([]entities.FlashCardDTO, len(words))
	for k := range words {
		dto[k].Translation = words[k].Translation
		dto[k].Word = words[k].Word
		dto[k].Lang = words[k].Lang
	}

	return dto
}
