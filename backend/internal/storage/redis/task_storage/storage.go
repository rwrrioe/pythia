package redis_storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
)

type RedisTaskStorage struct {
	ttl    time.Duration
	client *redis.Client
}

type TaskDTO struct {
	OCRText    []string                `json:"ocr_text"`
	Words      []entities.UnknownWord  `json:"words"`
	Examples   []entities.Example      `json:"examples"`
	FlashCards []entities.FlashCardDTO `json:"flashcards"`
}

func NewRedisTaskStorage(add string, ttl time.Duration) *RedisTaskStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     add,
		Password: "",
		DB:       0,
	})

	return &RedisTaskStorage{
		ttl:    ttl,
		client: rdb,
	}
}

func (s *RedisTaskStorage) Save(ctx context.Context, taskID string, task *TaskDTO) error {
	b, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return s.client.WithContext(ctx).Set("task:"+taskID, b, s.ttl).Err()
}

func (s *RedisTaskStorage) Get(ctx context.Context, taskID string) (*TaskDTO, bool, error) {
	val, err := s.client.WithContext(ctx).Get("task:" + taskID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, true, err
	}

	var task TaskDTO
	if err := json.Unmarshal([]byte(val), &task); err != nil {
		return nil, true, err
	}
	return &task, true, nil
}

func (s *RedisTaskStorage) UpdateTask(ctx context.Context, taskID string, update func(task *TaskDTO)) (bool, error) {
	val, err := s.client.WithContext(ctx).Get("task:" + taskID).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	var task TaskDTO
	if err := json.Unmarshal([]byte(val), &task); err != nil {
		return true, err
	}

	update(&task)
	b, err := json.Marshal(task)
	if err != nil {
		return true, err
	}

	err = s.client.WithContext(ctx).Set("task:"+taskID, b, time.Hour).Err()
	if err != nil {
		return true, err
	}

	return true, nil
}

func (s *RedisTaskStorage) Delete(ctx context.Context, taskID string) error {
	return s.client.WithContext(ctx).Del("task:" + taskID).Err()
}
