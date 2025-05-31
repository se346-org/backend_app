package redis

import (
	"context"
	"fmt"

	"github.com/chat-socio/backend/configuration"
	"github.com/redis/go-redis/v9"
)

func Connect(redisConf *configuration.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: redisConf.Password,
		DB:       redisConf.Database,
	})

	if err := client.Ping(context.TODO()).Err(); err != nil {
		panic(err)
	}

	return client
}
