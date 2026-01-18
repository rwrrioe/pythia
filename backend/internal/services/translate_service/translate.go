package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/domain/models"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
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

const defaultPrompt string = `
Ты профессиональный переводчик. 
Определи сложные или неизвестные слова в тексте на основе уровня "%s" и длительности изучения "%s".
Родной язык учащегося - русский. Выбери только 10-15 самых сложных или вероятно непонятных слов. Поставь слова в именительный падеж, настоящее время(как в словаре).
Дай перевод на русский в формате JSON [{"word": "...", "translation": "..."}].
Текст: %s`

const examplePrompt string = `
Ты профессиональный переводчик. 
Ниже дан список найденных незнакомых слов и исходный текст. 
Сделай контекстные примеры использования этих слов на основе уровня "%s" и длительности изучения "%s".
Дай перевод на русский в формате JSON [{"word": "...", "example": "..."}].
Текст: %s, слова %s`

func (t *TranslateService) FindUnknownWords(ctx context.Context, task *taskstorage.TaskDTO, req models.AnalyzeRequest) ([]entities.UnknownWord, error) {
	if task.OCRText == nil {
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
	txt := strings.Join(task.OCRText, " ")
	prompt := fmt.Sprintf(defaultPrompt, req.Level, req.Durating, txt)

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

func (t *TranslateService) WriteExamples(ctx context.Context, task *taskstorage.TaskDTO, req models.AnalyzeRequest) ([]entities.Example, error) {
	if task.OCRText == nil {
		return nil, errors.New("empty text in request")
	}
	var examples []entities.Example

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"word":    {Type: genai.TypeString},
					"example": {Type: genai.TypeString},
				},
				Required: []string{"word", "example"},
			},
		},
	}

	txt := strings.Join(task.OCRText, " ")
	b, err := json.Marshal(task.Words)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(examplePrompt, req.Level, req.Durating, txt, string(b))
	result, err := t.client.Models.GenerateContent(ctx,
		t.model,
		genai.Text(prompt),
		config,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate AI examples-response:%w", err)
	}

	if err := json.Unmarshal([]byte(result.Text()), &examples); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AI examples-response: %w", err)
	}

	return examples, nil
}
