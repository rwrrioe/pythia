package redis_storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
)

type RedisStorage struct {
	ttl    time.Duration
	client *redis.Client
}

type SessionDTO struct {
	Id        int64           `json:"session_id"`
	Name      string          `json:"name"`
	UserId    int64           `json:"user_id"`
	StartedAt time.Time       `json:"started_at"`
	EndedAt   time.Time       `json:"ended_at"`
	Duration  time.Duration   `json:"duration"`
	Status    string          `json:"status"`
	Language  int             `json:"language"`
	Level     int             `json:"level"`
	Accuracy  float64         `json:"accuracy"`
	Words     []entities.Word `json:"imp_words"`
}
type TaskDTO struct {
	SessionId int64           `json:"session_id"`
	OCRText   []string        `json:"ocr_text"`
	Words     []entities.Word `json:"words"`
}

func NewRedisStorage(ctx context.Context, add string, ttl time.Duration) (*RedisStorage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     add,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	//creating index
	_, err := rdb.FTCreate(
		ctx,
		"idx:tasks",
		// Options:
		&redis.FTCreateOptions{
			OnJSON: true,
			Prefix: []interface{}{"task:"},
		},
		// Index schema fields:
		&redis.FieldSchema{
			FieldName: "$.session_id",
			As:        "sessionId",
			FieldType: redis.SearchFieldTypeNumeric,
		},
	).Result()

	if err != nil && !strings.Contains(err.Error(), "Index already exists") {
		return nil, err
	}
	return &RedisStorage{
		ttl:    ttl,
		client: rdb,
	}, nil
}

func (s *RedisStorage) Save(ctx context.Context, taskId string, task TaskDTO) error {
	key := fmt.Sprintf("task:%s", taskId)
	if err := s.client.JSONSet(ctx, key, "$", task).Err(); err != nil {
		return err
	}

	_ = s.client.Expire(ctx, key, s.ttl).Err()
	return nil
}

func (s *RedisStorage) Get(ctx context.Context, taskId string) (*TaskDTO, bool, error) {
	key := fmt.Sprintf("task:%s", taskId)

	val, err := s.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, true, err
	}

	var arr []TaskDTO
	if err := json.Unmarshal([]byte(val), &arr); err != nil {
		return nil, true, err
	}
	if len(arr) == 0 {
		return nil, false, nil
	}

	return &arr[0], true, nil
}

func (s *RedisStorage) GetBySession(ctx context.Context, sessionId int64) ([]TaskDTO, bool, error) {
	q := fmt.Sprintf("@sessionId:[%d %d]", sessionId, sessionId)

	res, err := s.client.FTSearchWithArgs(ctx, "idx:tasks", q, &redis.FTSearchOptions{
		DialectVersion: 2,
		Return:         []redis.FTSearchReturn{{FieldName: "$"}},
	}).Result()

	if err != nil {
		return nil, false, err
	}

	if res.Total == 0 {
		return nil, false, nil
	}

	tasks := make([]TaskDTO, 0, len(res.Docs))
	for _, doc := range res.Docs {
		raw, ok := doc.Fields["$"]
		if !ok {
			continue
		}
		var t TaskDTO

		if err := json.Unmarshal([]byte(raw), &t); err != nil {
			return nil, true, err
		}

		tasks = append(tasks, t)
	}

	return tasks, true, nil
}

func (s *RedisStorage) UpdateTask(ctx context.Context, taskId string, update func(task *TaskDTO)) (bool, error) {
	key := fmt.Sprintf("task:%s", taskId)

	val, err := s.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	var arr []TaskDTO
	if err := json.Unmarshal([]byte(val), &arr); err != nil {
		return true, err
	}
	task := arr[0]

	update(&task)

	err = s.client.JSONSet(ctx, key, "$", task).Err()
	if err != nil {
		return true, err
	}

	_ = s.client.Expire(ctx, key, s.ttl)
	return true, nil
}

func (s *RedisStorage) Delete(ctx context.Context, taskID string) error {
	return s.client.Del(ctx, "task:"+taskID).Err()
}

func (s *RedisStorage) SaveSession(ctx context.Context, ss SessionDTO) error {
	key := fmt.Sprintf("session:%d", ss.Id)

	b, err := json.Marshal(ss)
	if err != nil {
		return err
	}
	if err := s.client.Set(ctx, key, b, s.ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (s *RedisStorage) GetSession(ctx context.Context, ssId int64) (*SessionDTO, bool, error) {
	key := fmt.Sprintf("session:%d", ssId)

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, true, err
	}

	var ss SessionDTO
	if err := json.Unmarshal([]byte(val), &ss); err != nil {
		return nil, true, err
	}
	return &ss, true, nil
}

func (s *RedisStorage) UpdateSession(ctx context.Context, ssId int64, update func(s *SessionDTO)) (bool, error) {
	key := fmt.Sprintf("session:%d", ssId)

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	var ss SessionDTO
	if err := json.Unmarshal([]byte(val), &ss); err != nil {
		return true, err
	}

	update(&ss)
	b, err := json.Marshal(ss)
	if err != nil {
		return true, err
	}

	err = s.client.Set(ctx, key, b, s.ttl).Err()
	if err != nil {
		return true, err
	}

	return true, nil
}
