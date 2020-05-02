package sources

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/lissein/dsptch/shared"
)

type RedisSource struct {
	client *redis.Client
	config *SourceConfig
}

func NewRedisSource(config *SourceConfig) (*RedisSource, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	source := &RedisSource{
		client,
		config,
	}
	return source, nil
}

func (source *RedisSource) Listen(messages chan shared.SourceMessage) error {
	pubsub := source.client.Subscribe("test", "blah")

	for {
		message := <-pubsub.Channel()
		messages <- shared.SourceMessage{
			Source:  fmt.Sprintf("redis/%s", message.Channel),
			Content: message.Payload,
		}
	}
}
