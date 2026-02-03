package service

import (
	"context"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
)

type FlashCardProvider interface {
	List(ctx context.Context, q postgresql.Querier, uid int64) ([]entities.FlashCard, error)
	ListByDeck(ctx context.Context, q postgresql.Querier, deckId int, uid int64) ([]entities.FlashCard, error)
	GetOrCreate(ctx context.Context, q postgresql.Querier, flCard entities.FlashCard, uid int64) (int, error)
}

type DeckProvider interface {
	ListBySession(ctx context.Context, q postgresql.Querier, sessionId int64, uid int64) (*entities.Deck, error)
	AttachFlashcard(ctx context.Context, q postgresql.Querier, deckId int, flId int) error
	GetOrCreate(ctx context.Context, q postgresql.Querier, sessionId int64, uid int64) (int, error)
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
