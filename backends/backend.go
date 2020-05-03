package backends

import "go.uber.org/zap"

type Constructor func(*Config) (Backend, error)

type InitFunction = func() (Backend, error)

type Config struct {
	Logger *zap.SugaredLogger
	Config map[string]interface{}
}

type Message struct {
	Source  string
	Payload interface{}
}

type Backend interface {
	Listen(messages chan Message)
	Handle(message Message) error
}
