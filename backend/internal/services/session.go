package service

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/auth/authz"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/domain/requests"
	service "github.com/rwrrioe/pythia/backend/internal/services/errors"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
)

type SessionProvider interface {
	GetSession(ctx context.Context, q postgresql.Querier, sessionId int64, uid int64) (*entities.Session, error)
	ListSessions(ctx context.Context, q postgresql.Querier, uid int64) ([]entities.Session, error)
	ListLatest(ctx context.Context, q postgresql.Querier, uid int64) ([]entities.Session, error)
	SaveSession(ctx context.Context, q postgresql.Querier, ss entities.Session, uid int64) (int, error)
	TryMarkFinished(ctx context.Context, q postgresql.Querier, sessionId int64, uid int64, endedAt time.Time) (bool, error)
	UpdateAccuracy(ctx context.Context, q postgresql.Querier, sessionId int64, uid int64, accuracy float64) error
}

const (
	Finished = "finished"
	Active   = "active"
)

type SessionService struct {
	OCR        *OCRService
	Translate  *TranslateService
	Learn      *LearnService
	Flashcards *FlashCardsService
	Redis      *taskstorage.RedisStorage

	txm *postgresql.TxManager

	SessionProvider    SessionProvider
	DeckProvider       DeckProvider
	FlashCardsProvider FlashCardProvider
	authorizer         authz.AuthorizeService
}

func NewSessionService(
	ocr *OCRService,
	transl *TranslateService,
	learn *LearnService,
	fl *FlashCardsService,
	redis *taskstorage.RedisStorage,
	txm *postgresql.TxManager,
	ss SessionProvider,
	deck DeckProvider,
	flProvider FlashCardProvider,
	authz authz.AuthorizeService,
) (*SessionService, error) {

	return &SessionService{
		OCR:                ocr,
		Translate:          transl,
		Learn:              learn,
		Redis:              redis,
		txm:                txm,
		SessionProvider:    ss,
		DeckProvider:       deck,
		FlashCardsProvider: flProvider,
		Flashcards:         fl,
		authorizer:         authz,
	}, nil
}

// sesison userflow
func (s *SessionService) StartSession(ctx context.Context, req requests.CreateSession) (int64, error) {
	const op = "service.SessionService.NewSession"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return 0, fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
	}

	ssion := entities.Session{
		UserId:    uid,
		Duration:  time.Duration(req.Duration) * time.Second,
		Status:    Active,
		Language:  req.LangId,
		StartedAt: time.Now(),
	}

	sessionId, err := s.SessionProvider.SaveSession(ctx, s.txm.Pool, ssion, uid)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	if err := s.Redis.SaveSession(ctx, taskstorage.SessionDTO{
		Id:        int64(sessionId),
		UserId:    ssion.UserId,
		Name:      ssion.Name,
		Status:    ssion.Status,
		Language:  ssion.Language,
		StartedAt: ssion.StartedAt,
	}); err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return int64(sessionId), nil
}

func (s *SessionService) RecognizeText(ctx context.Context, sessionId int64, taskId string, data []byte, lang string) error {
	const op = "service.SessionService.RecognizeText"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
	}
	if err := s.authorizer.CanAccessSession(ctx, uid, sessionId); errors.Is(err, authz.ErrForbidden) {
		return fmt.Errorf("%s:%w", op, service.ErrForbidden)
	}

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

func (s *SessionService) FindWords(ctx context.Context, sessionId int64, taskId string) ([]entities.Word, error) {
	const op = "service.SessionService.FindWords"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
	}
	if err := s.authorizer.CanAccessSession(ctx, uid, sessionId); errors.Is(err, authz.ErrForbidden) {
		return nil, fmt.Errorf("%s:%w", op, service.ErrForbidden)
	}

	t, ok, err := s.Redis.Get(ctx, taskId)
	if ok != true {
		return nil, fmt.Errorf("%s:%s", op, service.ErrTaskNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if ok != true {
		return nil, fmt.Errorf("%s:%s", op, service.ErrSessionNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	words, err := s.Translate.FindUnknownWords(ctx, t, requests.AnalyzeRequest{
		Level: levelsMap[ss.Level],
		Lang:  langsMap[ss.Language],
	})

	if ok, err = s.Redis.UpdateTask(ctx, taskId, func(task *taskstorage.TaskDTO) {
		task.Words = words
	}); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	} else if ok != true {
		return nil, fmt.Errorf("%s:%s", op, service.ErrTaskNotFound)
	}

	return words, nil
}

func (s *SessionService) EndSession(ctx context.Context, sessionId int64) error {
	const op = "service.SessionService.EndSession"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
	}
	if err := s.authorizer.CanAccessSession(ctx, uid, sessionId); errors.Is(err, authz.ErrForbidden) {
		return fmt.Errorf("%s:%w", op, service.ErrForbidden)
	}

	//find the most important words read redis + call translate
	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	if ok != true {
		return fmt.Errorf("%s:%s", op, service.ErrSessionNotFound)
	}

	tasks, ok, err := s.Redis.GetBySession(ctx, sessionId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%s:%w", op, service.ErrNoWords)
		}

		return fmt.Errorf("%s:%w", op, err)
	}
	if !ok {
		return fmt.Errorf("%s:%w", op, service.ErrSessionNotFound)
	}

	var words []entities.Word
	for _, t := range tasks {
		ws := t.Words
		for _, w := range ws {
			words = append(words, w)
		}
	}

	impWords, err := s.Translate.SummarizeWords(ctx, words, requests.AnalyzeRequest{
		Level: levelsMap[ss.Level],
		Lang:  langsMap[ss.Language],
	})
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	endedAt := time.Now()

	//save to the db
	err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		ok, err := s.SessionProvider.TryMarkFinished(ctx, tx, sessionId, uid, endedAt)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		deckId, err := s.DeckProvider.GetOrCreate(ctx, tx, sessionId, uid)
		if err != nil {
			return err
		}

		for _, w := range impWords {
			flId, err := s.FlashCardsProvider.GetOrCreate(ctx, tx, entities.FlashCard{
				Word:   w.Word,
				Transl: w.Translation,
				Lang:   extractLangId(w.Lang),
			}, uid)

			if err != nil {
				return err
			}

			if err := s.DeckProvider.AttachFlashcard(ctx, tx, deckId, flId); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	//commit + redis update
	_, _ = s.Redis.UpdateSession(ctx, sessionId, func(dto *taskstorage.SessionDTO) {
		dto.Status = Finished
		dto.EndedAt = endedAt
		dto.Words = impWords
	})

	return nil
}

func (s *SessionService) GetFlashcards(ctx context.Context, sessionId int64) ([]entities.FlashCardDTO, error) {
	const op = "service.SessionService.GetFlashcards"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
	}
	if err := s.authorizer.CanAccessSession(ctx, uid, sessionId); errors.Is(err, authz.ErrForbidden) {
		return nil, fmt.Errorf("%s:%w", op, service.ErrForbidden)
	}

	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if ok != true {
		return nil, fmt.Errorf("%s:%s", op, service.ErrSessionNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	flCards := s.Flashcards.BuildCards(ctx, ss.Words)
	return flCards, nil
}

func (s *SessionService) Quiz(ctx context.Context, sessionId int64) ([]entities.QuizQuestion, error) {
	const op = "service.SessionService.Quiz"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
	}
	if err := s.authorizer.CanAccessSession(ctx, uid, sessionId); errors.Is(err, authz.ErrForbidden) {
		return nil, fmt.Errorf("%s:%w", op, service.ErrForbidden)
	}

	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if ok != true {
		return nil, fmt.Errorf("%s:%s", op, service.ErrSessionNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	quiz := s.Learn.QuizTest(ctx, ss.Words)
	return quiz, nil
}

func (s *SessionService) SummarizeSession(ctx context.Context, sessionId int64, accuracy float64) ([]entities.Word, error) {
	const op = "service.SessionService.SummarizeSession"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
	}
	if err := s.authorizer.CanAccessSession(ctx, uid, sessionId); errors.Is(err, authz.ErrForbidden) {
		return nil, fmt.Errorf("%s:%w", op, service.ErrForbidden)
	}

	if err := s.SessionProvider.UpdateAccuracy(ctx, s.txm.Pool, sessionId, uid, accuracy); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	ss, ok, err := s.Redis.GetSession(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	if ok != true {
		return nil, fmt.Errorf("%s:%s", op, service.ErrSessionNotFound)
	}

	return ss.Words, nil
}

var levelsMap = map[int]string{
	1: "A2",
	2: "B1",
	3: "B2",
}

var langsMap = map[int]string{
	1: "de",
	2: "en",
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
