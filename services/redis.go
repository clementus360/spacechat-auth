package services

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Connect to a redis instance using a URI from environment variables
func ConnectRedis() *redis.Client {
	REDIS_URI := os.Getenv("REDIS_URI")
	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_URI,
		Password: "",
		DB:       0,
	})

	return rdb
}

func StoreSession(userId string, SessionId string) error {
	rdb := ConnectRedis()

	err := rdb.Set(context.Background(), userId, SessionId, time.Minute*10).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetSession(userId string) (string, error) {
	rdb := ConnectRedis()

	Session, err := rdb.Get(context.Background(), Hash(userId)).Result()
	if err != nil {
		return "", err
	}

	return Session, nil
}
