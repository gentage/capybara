package pubsub

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

type RedisClient struct {
	client *redis.Client
}

func MakeRedisClient(host string, port int) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(
			&redis.Options{
				Addr: fmt.Sprintf("%s:%d", host, port),
			},
		),
	}
}

func (r *RedisClient) Publish(channel string, msg string) error {
	r.client.Publish(channel, msg)
	return nil
}

func (r *RedisClient) Subscribe(channel string) chan string {
	c := make(chan string)
	messages := r.client.Subscribe(channel).ChannelSize(16)
	go func() {
		for m := range messages {
			c <- m.Payload
		}
	}()
	return c
}
