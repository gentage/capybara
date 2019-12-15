package pubsub

type Client interface {
	Publish(string, string) error
	Subscribe(string) (<-chan string, error)
}
