package redis

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

type QueueRepository interface {
	Publish() (string, error)
	Read()
}

func (r *RedisQueueRepository) Publish() (string, error) {
	args := redis.XAddArgs{
		Stream: r.Stream,
		Values: map[string]interface{}{
			"action":   "start",
			"resource": "product-1",
			"user":     "bob",
		},
	}
	return r.Client.XAdd(&args).Result()
}

func Process(r *RedisQueueRepository) {
	id := "0"
	for {
		args := redis.XReadArgs{
			Streams: []string{r.Stream, id},
			Count:   1,
			Block:   0,
		}
		data, err := r.Client.XRead(&args).Result()
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
	Client *redis.Client
	Stream string
}

func NewQueue(client *redis.Client, stream string) (*RedisQueueRepository, error) {
	rp := RedisQueueRepository{Client: client, Stream: stream}
	rp.Read()
	return &rp, nil
}
