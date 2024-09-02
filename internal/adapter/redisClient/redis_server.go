package RedisCache

import (
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	Client *redis.ClusterClient
}

func NewRedisRepository(client *redis.ClusterClient) *RedisRepository {
	return &RedisRepository{
		Client: client,
	}
}
