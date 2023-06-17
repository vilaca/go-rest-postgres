package redis

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

type QueueRepository interface {
	// CreateUser(user *models.User) (*models.User, error)
	// GetUser(id string) (*models.User, error)
	// UserExists(name string) (bool, error)
	// Login(name string, password string) (bool, error)
	// CreateSession(*models.Session) error
	Publish() (string, error)
	Read()
}

func (r *RedisQueueRepository) Publish() (string, error) {
	args := redis.XAddArgs{
		Stream: "ohoh-stream",
		Values: map[string]interface{}{
			"action":   "start",
			"resource": "product-1",
			"user":     "bob",
		},
	}
	return r.client.XAdd(&args).Result()
}

func Process(r *RedisQueueRepository) {
	id := "0"
	for {
		args := redis.XReadArgs{
			Streams: []string{"ohoh-stream", id},
			Count:   1,
			Block:   0,
		}
		data, err := r.client.XRead(&args).Result()
		if err != nil {
			log.Fatal(err)
		}
		for _, result := range data {
			for _, message := range result.Messages {
				fmt.Println("Job received is", message)
				id = message.ID
			}
		}
	}
}

func (r *RedisQueueRepository) Read() {
	go Process(r)
}

type RedisQueueRepository struct {
	client *redis.Client
}

func NewQueue(client *redis.Client) (*RedisQueueRepository, error) {
	return &RedisQueueRepository{client: client}, nil
}
