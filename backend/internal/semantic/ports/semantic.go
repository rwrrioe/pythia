package ports

import (
	"context"

	"github.com/rwrrioe/pythia/backend/internal/semantic/domain"
)

type Translator interface {
	TranslateUnknown(
		ctx context.Context,
		text string,
		opts *Options,
	) ([]domain.Word, error)

	FindMostImportant(
		ctx context.Context,
		words []domain.Word,
		opts *Options,
	) ([]domain.Word, error)
}

type Options struct {
	Level    string
	Durating string
	Language string
}
