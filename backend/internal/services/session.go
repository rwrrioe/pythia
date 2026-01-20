package service

import (
	"context"

	"fmt"
	"math/rand"
	"time"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/domain/requests"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
)

type SessionProvider interface {
	GetSession(ctx context.Context, sessionId int) (*entities.Session, error)
	ListSessions(ctx context.Context) ([]entities.Session, error)
	ListLatest(ctx context.Context) ([]entities.Session, error)
}

type SessionService struct {
	OCR                *OCRService
	Translate          *TranslateService
	Learn              *LearnService
	Redis              *taskstorage.RedisStorage
	SessionProvider    SessionProvider
	DeckProvider       DeckProvider
	FlashCardsProvider FlashCardProvider
	Session            entities.Session
}

func (s *SessionService) NewSession(
	ctx context.Context,
	req requests.CreateSession,
	ocr *OCRService,
	transl *TranslateService,
	learn *LearnService,
	redis *taskstorage.RedisStorage,
	ss SessionProvider,
	deck DeckProvider,
	fl FlashCardProvider,
) (*SessionService, error) {
	const op = "service.SessionService.NewSession"
	//uid, ok := auth.UIDFromContext(ctx)
	//if !ok {
	//	return nil, fmt.Errorf("%s:%w", op, errors.New("user not found"))
	//}

	ssion := entities.Session{
		Id:        rand.Int(),
		Duration:  req.Duration,
		Status:    "active",
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
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &SessionService{
		OCR:                ocr,
		Translate:          transl,
		Learn:              learn,
		Redis:              redis,
		SessionProvider:    ss,
		DeckProvider:       deck,
		FlashCardsProvider: fl,
		Session:            ssion,
	}, nil
}

func (s *SessionService) EndSession(ctx context.Context, sessionId int, topic string) error {
	const op = "service.SessionService.EndSession"

	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if ok != true {
		return fmt.Errorf("%s:%s", op, "session not found")
	}
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	ok, err = s.Redis.UpdateSession(ctx, ss.Id, func(ss *taskstorage.SessionDTO) {
		ss.EndedAt = time.Now()
		ss.Name = topic
	})

	if ok != true {
		return fmt.Errorf("%s:%s", op, "session not found")
	}

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *SessionService) RecognizeText(ctx context.Context, taskId string, data []byte, lang string) error {
	const op = "service.SessionService.RecognizeText"

	txt, err := s.OCR.ProcessImage(ctx, data, lang)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)

	}

	if err := s.Redis.Save(ctx, taskId, taskstorage.TaskDTO{
		SessionId: s.Session.Id,
		OCRText:   txt,
	}); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (s *SessionService) FindWords(ctx context.Context, taskId string) error {
	const op = "service.SessionService.FindWords"

	t, ok, err := s.Redis.Get(ctx, taskId)
	if ok != true {
		return fmt.Errorf("%s:%s", op, "task not found")
	}

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	words, err := s.Translate.FindUnknownWords(ctx, t, requests.AnalyzeRequest{
		Level: levelsMap[s.Session.Level],
		Lang:  langsMap[s.Session.Language],
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

	words, err := s.Translate.SummarizeWords(ctx, sessionId, t, requests.AnalyzeRequest{
		Level: levelsMap[s.Session.Level],
		Lang:  langsMap[s.Session.Language],
	})
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return words, err
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
