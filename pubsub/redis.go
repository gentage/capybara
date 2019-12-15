package pubsub

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(host string, port int) Client {
	return &redisClient{
		client: redis.NewClient(
			&redis.Options{
				Addr: fmt.Sprintf("%s:%d", host, port),
			},
		),
	}
}

func (r *redisClient) Publish(channel string, msg string) error {
	r.client.Publish(channel, msg)
	return nil
}

func (r *redisClient) Subscribe(channel string) (<-chan string, error) {
	c := make(chan string)
	messages := r.client.Subscribe(channel).ChannelSize(16)
	go func() {
		for m := range messages {
			c <- m.Payload
		}
	}()
	return c, nil
}
