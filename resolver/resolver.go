package resolver

import (
	"context"

	"github.com/gentage/capybara/pubsub"
)

type Resolver struct {
	pubsubClient pubsub.Client
}

func MakeResolver(pubsubClient pubsub.Client) *Resolver {
	return &Resolver{pubsubClient: pubsubClient}
}

func (r *Resolver) Ping() string {
	return "Pong!"
}

type PublishArgs struct {
	Channel string
	Msg     string
}

func (r *Resolver) Publish(args PublishArgs) string {
	_ = r.pubsubClient.Publish(args.Channel, args.Msg)
	return args.Msg
}

type SubscribeArgs struct {
	Channel string
}

func (r *Resolver) Subscribe(ctx context.Context, args SubscribeArgs) <-chan string {
	return r.pubsubClient.Subscribe(args.Channel)
}
