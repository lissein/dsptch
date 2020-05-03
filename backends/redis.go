package backends

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

type RedisBackend struct {
	client *redis.Client
	config *BackendConfig

	channels []string
}

func NewRedisBackend(config *BackendConfig) *RedisBackend {
	backend := &RedisBackend{
		config:   config,
		channels: config.Config["channels"].([]string),
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}

	_, err := backend.client.Ping().Result()
	if err != nil {
		config.Logger.Panic(err)
	}

	return backend
}

func (backend *RedisBackend) Listen(messages chan BackendInputMessage) {
	pubsub := backend.client.Subscribe(backend.channels...)

	for {
		message := <-pubsub.Channel()
		messages <- BackendInputMessage{
			Source:  fmt.Sprintf("redis/%s", message.Channel),
			Content: message.Payload,
		}
	}
}

func (backend *RedisBackend) HandleMessage(message BackendOutputMessage) error {
	return nil
}
