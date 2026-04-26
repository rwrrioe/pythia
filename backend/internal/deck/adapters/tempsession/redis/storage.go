package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	TTL time.Duration
}

type Storage struct {
	client *redis.Client
	config *Config
}

func New(
	cl *redis.Client,
	cfg *Config,
) *Storage {
	return &Storage{
		client: cl,
		config: cfg,
	}
}
