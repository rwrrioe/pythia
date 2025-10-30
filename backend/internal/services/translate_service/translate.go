package translate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rwrrioe/pythia/internal/domain/entities"
	"github.com/rwrrioe/pythia/internal/domain/models"
	"google.golang.org/genai"
)

type TranslateService struct {
	client genai.Client
	model  string
}

func NewTranslateService(ctx context.Context, model string) (*TranslateService, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &TranslateService{client: *client, model: model}, nil
}

const defaultPrompt string = `Определи сложные или неизвестные слова в тексте на основе уровня "%s" и длительности изучения "%s".
Дай перевод на русский в формате JSON [{"word": "...", "translation": "..."}].
Текст: %s`

func (t *TranslateService) FindUnknownWords(ctx context.Context, req models.AnalyzeRequest) ([]entities.UnknownWord, error) {
	if string(req.Text) == "" {
		return nil, errors.New("empty text in request")
	}
	var words []entities.UnknownWord

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"word":        {Type: genai.TypeString},
					"translation": {Type: genai.TypeString},
				},
				Required: []string{"word", "translation"},
			},
		},
	}

	prompt := fmt.Sprintf(defaultPrompt, req.Level, req.Durating, req.Text)

	result, err := t.client.Models.GenerateContent(ctx,
		t.model,
		genai.Text(prompt),
		config,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response:%w", err)
	}

	if err := json.Unmarshal([]byte(result.Text()), &words); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AI response: %w", err)
	}

	for i := range words {
		words[i].Lang = req.Lang
	}

	return words, nil
}
