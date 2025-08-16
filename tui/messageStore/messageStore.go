package messageStore

import (
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

func Init() {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		ctx := context.Background()

		_, err := client.Ping(ctx).Result()
		if err != nil {
			log.Fatal(err)
		}
	})
}

func GetClient() *redis.Client {
	if client == nil {
		Init()
	}
	return client
}
