package translate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

type TranslateService struct{}

func NewTranslateService() *TranslateService {
	return &TranslateService{}
}

type UnknownWord struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

// find unknown words from the text
func (t *TranslateService) FindUnknown(ctx context.Context, text string) ([]UnknownWord, error) {
	var words []UnknownWord
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load client: %w", err)
	}

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

	prompt := fmt.Sprint(`
	Ты - ИИ агент, созданный для помощи изучающим немецкий и английский язык на уровнях A2 - B2. Пользователь отправляет текст с учебника, в понимании которого он испытывает трудности.
	Ты должен найти возможно неизвестные пользователю слова на основе его уровня, учебника, длительности изучения языка. После того как ты нашел неизвестное слово, ты должен перевести его на русский язык с сохранением контекста.
	Строго соблюдай структуру и не добавляй лишнего.
	`, text)

	result, err := client.Models.GenerateContent(ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response:%w", err)
	}

	if err := json.Unmarshal([]byte(result.Text()), &words); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AI response: %w", err)
	}

	return words, nil
}
