package service

import (
	"context"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
)

type FlashCardsService struct{}

func NewCardsService() *FlashCardsService {
	return &FlashCardsService{}
}

func (c *FlashCardsService) BuildCards(ctx context.Context, words []entities.UnknownWord) ([]entities.FlashCardDTO, error) {
	dto := make([]entities.FlashCardDTO, len(words))
	for k := range words {
		dto[k].Translation = words[k].Translation
		dto[k].Word = words[k].Word
		dto[k].Lang = words[k].Lang
	}

	return dto, nil
}
