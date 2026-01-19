package service

import (
	"context"
	"math/rand"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
)

type LearnService struct {
	Optcount int
}

func NewLearnService(optcount int) *LearnService {
	return &LearnService{Optcount: optcount}
}

func (s *LearnService) QuizTest(ctx context.Context, words *[]entities.Word) []entities.QuizQuestion {
	var test []entities.QuizQuestion

	for _, v := range *words {
		opts := pickOptions(words, v.Word, s.Optcount)
		questionDTO := entities.QuizQuestion{
			Answer:   v.Word,
			Question: v.Translation,
			Options:  opts,
		}
		test = append(test, questionDTO)
	}

	return test

}

func pickOptions(words *[]entities.Word, correct string, optcount int) []string {
	var pool []string

	for _, w := range *words {
		if w.Word != correct {
			pool = append(pool, w.Word)
		}
	}

	rand.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	opts := make([]string, 0, optcount)
	opts = append(opts, pool[:4]...)

	insIdx := rand.Intn(len(opts))
	opts[insIdx] = correct

	return opts
}
