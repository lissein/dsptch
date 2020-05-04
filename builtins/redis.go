package builtins

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/lissein/dsptch/backends"
)

type RedisBackend struct {
	client *redis.Client
	config *backends.Config

	channels []string
}

func NewRedisBackend(config *backends.Config) (backends.Backend, error) {
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

	return backend, nil
}

func (backend *RedisBackend) Listen(messages chan backends.Message) {
	pubsub := backend.client.Subscribe(backend.channels...)

	for {
		message := <-pubsub.Channel()
		messages <- backends.Message{
			Source:  fmt.Sprintf("redis/%s", message.Channel),
			Payload: message.Payload,
		}
	}
}

func (backend *RedisBackend) Handle(message backends.Message) error {
	return nil
}
