package backends

import "go.uber.org/zap"

type BackendConstructor func(*BackendConfig) Backend

type BackendConfig struct {
	Logger *zap.SugaredLogger
	Config map[string]interface{}
}

type BackendInputMessage struct {
	Source  string
	Content interface{}
}

type BackendOutputMessage struct {
	Content interface{}
}

type Backend interface {
	// Listen and publish messages to the `messages` channel so that the app can handle it
	Listen(messages chan BackendInputMessage)

	// Handle a message
	HandleMessage(message BackendOutputMessage) error
}
