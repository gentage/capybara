package resolver

import (
	"fmt"
	"time"

	"github.com/gentage/capybara/pubsub"
)

type resolver struct {
	pubsubClient pubsub.Client
}

func NewResolver(pubsubClient pubsub.Client) *resolver {
	return &resolver{pubsubClient: pubsubClient}
}

func (r *resolver) Ping() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type PublishArgs struct {
	Channel string
	Msg     string
}

func (r *resolver) Publish(args PublishArgs) string {
	_ = r.pubsubClient.Publish(args.Channel, args.Msg)
	return args.Msg
}

type SubscribeArgs struct {
	Channel string
}

func (r *resolver) Subscribe(args SubscribeArgs) <-chan string {
	c, _ := r.pubsubClient.Subscribe(args.Channel)
	return c
}
