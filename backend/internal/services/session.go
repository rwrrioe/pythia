package service

import (
	"context"
	"time"

	"fmt"
	"math/rand"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/domain/requests"
	service "github.com/rwrrioe/pythia/backend/internal/services/errors"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
)

type SessionProvider interface {
	GetSession(ctx context.Context, sessionId int) (*entities.Session, error)
	ListSessions(ctx context.Context) ([]entities.Session, error)
	ListLatest(ctx context.Context) ([]entities.Session, error)
	SaveSession(ctx context.Context, ss entities.Session) error
}

const (
	Finished = "finished"
	Active   = "active"
)

type SessionService struct {
	OCR                *OCRService
	Translate          *TranslateService
	Learn              *LearnService
	Flashcards         *FlashCardsService
	Redis              *taskstorage.RedisStorage
	SessionProvider    SessionProvider
	DeckProvider       DeckProvider
	FlashCardsProvider FlashCardProvider
}

func (s *SessionService) NewSessionService(
	ocr *OCRService,
	transl *TranslateService,
	learn *LearnService,
	fl *FlashCardsService,
	redis *taskstorage.RedisStorage,
	ss SessionProvider,
	deck DeckProvider,
	flProvider FlashCardProvider,
) (*SessionService, error) {
	const op = "service.SessionService.NewSession"
	//uid, ok := auth.UIDFromContext(ctx)
	//if !ok {
	//	return nil, fmt.Errorf("%s:%w", op, errors.New("user not found"))
	//}

	return &SessionService{
		OCR:                ocr,
		Translate:          transl,
		Learn:              learn,
		Redis:              redis,
		SessionProvider:    ss,
		DeckProvider:       deck,
		FlashCardsProvider: flProvider,
		Flashcards:         fl,
	}, nil
}

func (s *SessionService) StartSession(ctx context.Context, req requests.CreateSession) (int, error) {
	const op = "service.SessionService.NewSession"

	ssion := entities.Session{
		Id:        rand.Int(),
		Duration:  req.Duration,
		Status:    Active,
		Language:  req.LangId,
		StartedAt: time.Now(),
	}

	if err := s.Redis.SaveSession(ctx, taskstorage.SessionDTO{
		Id:        ssion.Id,
		Name:      ssion.Name,
		Status:    ssion.Status,
		Language:  ssion.Language,
		StartedAt: ssion.StartedAt,
	}); err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return ssion.Id, nil
}

func (s *SessionService) RecognizeText(ctx context.Context, sessionId int, taskId string, data []byte, lang string) error {
	const op = "service.SessionService.RecognizeText"

	txt, err := s.OCR.ProcessImage(ctx, data, lang)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)

	}

	if err := s.Redis.Save(ctx, taskId, taskstorage.TaskDTO{
		SessionId: sessionId,
		OCRText:   txt,
	}); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (s *SessionService) FindWords(ctx context.Context, sessionId int, taskId string) error {
	const op = "service.SessionService.FindWords"

	t, ok, err := s.Redis.Get(ctx, taskId)
	if ok != true {
		return fmt.Errorf("%s:%s", op, "task not found")
	}

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if ok != true {
		return fmt.Errorf("%s:%s", op, "session not found")
	}
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	words, err := s.Translate.FindUnknownWords(ctx, t, requests.AnalyzeRequest{
		Level: levelsMap[ss.Level],
		Lang:  langsMap[ss.Language],
	})

	if ok, err = s.Redis.UpdateTask(ctx, taskId, func(task *taskstorage.TaskDTO) {
		task.Words = words
	}); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	} else if ok != true {
		return fmt.Errorf("%s:%s", op, "task not found")
	}

	return nil
}

func (s *SessionService) SummarizeSession(ctx context.Context, sessionId int) ([]entities.Word, error) {
	const op = "service.SessionService.SummarizeSession"

	var words []entities.Word
	tasks, ok, err := s.Redis.GetBySession(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	} else if !ok {
		return nil, fmt.Errorf("%s:%w", op, service.ErrSessionNotFound)
	}

	for _, t := range tasks {
		ws := t.Words
		for _, w := range ws {
			words = append(words, w)
		}
	}
	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if ok != true {
		return nil, fmt.Errorf("%s:%s", op, "session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	impWords, err := s.Translate.SummarizeWords(ctx, words, requests.AnalyzeRequest{
		Level: levelsMap[ss.Level],
		Lang:  langsMap[ss.Language],
	})
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if _, err = s.Redis.UpdateSession(ctx, sessionId, func(s *taskstorage.SessionDTO) {
		s.Words = impWords
	}); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return impWords, err
}

func (s *SessionService) GetFlashcards(ctx context.Context, sessionId int) ([]entities.FlashCardDTO, error) {
	const op = "service.SessionService.GetFlashcards"

	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if ok != true {
		return nil, fmt.Errorf("%s:%s", op, "session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	flCards := s.Flashcards.BuildCards(ctx, ss.Words)
	return flCards, nil
}

func (s *SessionService) Quiz(ctx context.Context, sessionId int) ([]entities.QuizQuestion, error) {
	const op = "service.SessionService.Quiz"

	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if ok != true {
		return nil, fmt.Errorf("%s:%s", op, "session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	quiz := s.Learn.QuizTest(ctx, ss.Words)
	return quiz, nil
}

func (s *SessionService) EndSession(ctx context.Context, sessionId int, accuracy float64) error {
	const op = "service.SessionService.EndSession"

	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if ok != true {
		return fmt.Errorf("%s:%s", op, "session not found")
	}
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if ss.Status == Finished {
		return fmt.Errorf("%s:%w", op, service.ErrSessionAlreadyFinished)
	}

	deck := entities.Deck{
		Id:        rand.Int(),
		SessionId: sessionId,
	}
	if err := s.DeckProvider.SaveDeck(ctx, deck); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	for _, f := range ss.Words {
		fl := entities.FlashCard{
			Id:     rand.Int(),
			Word:   f.Word,
			Transl: f.Translation,
			Lang:   extractLangId(f.Lang),
		}

		if err := s.FlashCardsProvider.SaveFlashcard(ctx, fl); err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}

		if err := s.DeckProvider.AttachFlashcard(ctx, deck.Id, fl.Id); err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}
	}

	if err := s.SessionProvider.SaveSession(ctx, entities.Session{
		Id:        ss.Id,
		Name:      ss.Name,
		StartedAt: ss.StartedAt,
		Duration:  ss.Duration,
		Status:    Finished,
		Language:  ss.Language,
		Level:     ss.Level,
		Accuracy:  accuracy,
	}); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

var levelsMap = map[int]string{
	1: "A2",
	2: "B1",
	3: "B2",
}

var langsMap = map[int]string{
	1: "en",
	2: "de",
	3: "fr",
	4: "es",
}

func extractLangId(lang string) int {
	for k, v := range langsMap {
		if v == lang {
			return k
		}
	}
	return 0
}
