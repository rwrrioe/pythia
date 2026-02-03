package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/domain/requests"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
	"google.golang.org/genai"
)

type TranslateService struct {
	client genai.Client
	model  string
	Redis  *taskstorage.RedisStorage
}

func NewTranslateService(ctx context.Context, model string) (*TranslateService, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &TranslateService{client: *client, model: model}, nil
}

const findImportantPrompt string = `
You are a language learning expert and vocabulary curator.

You are given a list of words extracted from a single learning session.
The learner is at CEFR level %s.

Your task:
1. Analyze all the words together as a single session context.
2. Select ONLY 10–15 words that are the most important for active learning.
3. Prioritize words that:
   - are likely unknown or weakly known by an A2–B1 learner
   - are useful, high-value, or conceptually important
   - appear frequently or are central to the session topic
   - are not proper names or trivial function words
4. Deprioritize or exclude:
   - very basic words (A1 level)
   - words that are obvious from context or near-synonyms of simpler words
   - names, numbers, dates, or overly specific terms

Important:
- Think in terms of *learning value*, not raw frequency alone.
- The goal is efficient learning, not completeness.

Do NOT include any explanations outside the JSON.
Do NOT include more than 15 or fewer than 10 words.

Input words:
<<<
{{%s}}
>>>
`

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

func (t *TranslateService) FindUnknownWords(ctx context.Context, task *taskstorage.TaskDTO, req requests.AnalyzeRequest) ([]entities.Word, error) {
	if task.OCRText == nil {
		return nil, errors.New("empty text in request")
	}
	var words []entities.Word

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

func (t *TranslateService) WriteExamples(ctx context.Context, task *taskstorage.TaskDTO, req requests.AnalyzeRequest) ([]entities.Example, error) {
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

func (s *TranslateService) SummarizeWords(ctx context.Context, words []entities.Word, req requests.AnalyzeRequest) ([]entities.Word, error) {
	const op = "service.TranslateService.SummarizeWords"

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

	b, err := json.Marshal(words)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(findImportantPrompt, req.Level, string(b))
	result, err := s.client.Models.GenerateContent(ctx,
		s.model,
		genai.Text(prompt),
		config,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate AI examples-response:%w", err)
	}

	var found []entities.Word

	if err := json.Unmarshal([]byte(result.Text()), &found); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AI examples-response: %w", err)
	}

	for i := range found {
		found[i].Lang = req.Lang
	}
	return found, nil
}
