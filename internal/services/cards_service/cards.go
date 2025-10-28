package cards

import (
	"context"

	"github.com/rwrrioe/pythia/internal/domain/entities"
)

type CardsService struct{}

func NewCardsService() *CardsService {
	return &CardsService{}
}

func (c *CardsService) BuildCards(ctx context.Context, words []entities.UnknownWord) (*[]entities.FlashCardDTO, error) {
	dto := make([]entities.FlashCardDTO, len(words))
	for k := range words {
		dto[k].Translation = words[k].Translation
		dto[k].Word = words[k].Word
		dto[k].Lang = words[k].Lang
	}

	return &dto, nil
}
