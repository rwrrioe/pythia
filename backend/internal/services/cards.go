package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
)

type FlashCardProvider interface {
	List(ctx context.Context, q postgresql.Querier, uid int64) ([]entities.FlashCard, error)
	ListByDeck(ctx context.Context, q postgresql.Querier, deckId uuid.UUID, uid int64) ([]entities.FlashCard, error)
	GetOrCreate(ctx context.Context, q postgresql.Querier, flCard entities.FlashCard, uid int64) (uuid.UUID, error)
	FlashcardsPool() *pgxpool.Pool
}

type DeckProvider interface {
	ListBySession(ctx context.Context, q postgresql.Querier, sessionId uuid.UUID, uid int64) (*entities.Deck, error)
	AttachFlashcard(ctx context.Context, q postgresql.Querier, deckId uuid.UUID, flashcardId uuid.UUID) error
	GetOrCreate(ctx context.Context, q postgresql.Querier, sessionId uuid.UUID, uid int64) (uuid.UUID, error)
	DeckPool() *pgxpool.Pool
}

type FlashCardsService struct {
	flashcards FlashCardProvider
	decks      DeckProvider

	pool postgresql.Querier
}

func NewCardsService(
	flashcards FlashCardProvider,
	decks DeckProvider,
	pool postgresql.Querier,
) *FlashCardsService {
	return &FlashCardsService{
		flashcards: flashcards,
		decks:      decks,
		pool:       pool,
	}
}

func (s *FlashCardsService) BuildCards(ctx context.Context, words []entities.Word) []entities.FlashCardDTO {
	dto := make([]entities.FlashCardDTO, len(words))
	for k := range words {
		dto[k].Translation = words[k].Translation
		dto[k].Word = words[k].Word
		dto[k].Lang = words[k].Lang
	}

	return dto
}

func (s *FlashCardsService) GetBySession(ctx context.Context, sessionId uuid.UUID) ([]entities.FlashCard, error) {
	const op = "service.FlashcardService.GetBySession"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%s:%w", op, ErrUnauthorized)
	}

	deck, err := s.decks.ListBySession(ctx, s.pool, sessionId, uid)
	if err != nil {
		if errors.Is(err, postgresql.ErrDeckNotFound) {
			return nil, fmt.Errorf("%s:%w", op, ErrDeckNotFound)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}
	slog.Any(op, deck)

	flCards, err := s.flashcards.ListByDeck(ctx, s.pool, deck.Id, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	slog.Any(op, flCards)

	return flCards, nil
}
