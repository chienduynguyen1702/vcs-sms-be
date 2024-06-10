package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisProtocol, _ := strconv.Atoi(os.Getenv("REDIS_PROTOCOL"))

	redisOptions := &redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
		Protocol: redisProtocol,
	}

	rdb := redis.NewClient(redisOptions)

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
		return nil, err
	}

	log.Println("Connected to redis")
	return rdb, nil
}
