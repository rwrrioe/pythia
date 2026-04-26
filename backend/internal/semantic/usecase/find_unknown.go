package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/rwrrioe/pythia/backend/internal/lib/logger/sl"
	"github.com/rwrrioe/pythia/backend/internal/semantic/domain"
	"github.com/rwrrioe/pythia/backend/internal/semantic/ports"
)

type Semantic struct {
	semanticOperator ports.SemanticOperator
	log              *slog.Logger
}

func New(
	semanticOperator ports.SemanticOperator,
	log *slog.Logger,
) (*Semantic, error) {
	return &Semantic{
		semanticOperator: semanticOperator,
		log:              log,
	}, nil
}

func (s *Semantic) FindUnknownWords(
	ctx context.Context,
	text []string,
	opts *ports.Options,
) ([]domain.Word, error) {
	const op = "translate.FindUnknownWords"

	txt := strings.Join(text, " ")

	words, err := s.semanticOperator.TranslateUnknown(ctx, txt, opts)
	if err != nil {
		s.log.Error("failed to translate unknown words", sl.Err(err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	s.log.Info("unknown words successfully found")
	return words, nil
}

//func (t *Semantic) WriteExamples(ctx context.Context, task *taskstorage.TaskDTO, req requests.AnalyzeRequest) ([]domain.Example, error) {
//	if task.OCRText == nil {
//		return nil, errors.New("empty text in request")
//	}
//	var examples []domain.Example
//
//	config := &genai.GenerateContentConfig{
//		ResponseMIMEType: "application/json",
//		ResponseSchema: &genai.Schema{
//			Type: genai.TypeArray,
//			Items: &genai.Schema{
//				Type: genai.TypeObject,
//				Properties: map[string]*genai.Schema{
//					"word":    {Type: genai.TypeString},
//					"example": {Type: genai.TypeString},
//				},
//				Required: []string{"word", "example"},
//			},
//		},
//	}
//
//	txt := strings.Join(task.OCRText, " ")
//	b, err := json.Marshal(task.Words)
//	if err != nil {
//		return nil, err
//	}
//
//	prompt := fmt.Sprintf(examplePrompt, req.Level, req.Durating, txt, string(b))
//	result, err := t.client.Models.GenerateContent(ctx,
//		t.model,
//		genai.Text(prompt),
//		config,
//	)
//
//	if err != nil {
//		return nil, fmt.Errorf("failed to generate AI examples-response:%w", err)
//	}
//
//	if err := json.Unmarshal([]byte(result.Text()), &examples); err != nil {
//		return nil, fmt.Errorf("failed to unmarshal AI examples-response: %w", err)
//	}
//
//	return examples, nil
//}
