package backends

import (
	"github.com/go-redis/redis/v7"
)

type RedisBackend struct {
	client *redis.Client
	config *BackendConfig
}

func NewRedisBackend(config *BackendConfig) *RedisBackend {
	backend := &RedisBackend{
		config: config,
	}

	return backend
}

// Listen and publish messages to the `messages` channel so that the app can handle it
func (backend *RedisBackend) Listen(messages chan interface{}) {

}

// Handle a message
func (backend *RedisBackend) HandleMessage(message interface{}) error {
	return nil
}
