package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis:%v", err)

	}
	log.Println("Successfully connected to redis")
}
