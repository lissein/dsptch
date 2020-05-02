package backends

import (
	"fmt"
	"math/rand"
	"time"
)

type DummyBackend struct {
	config *BackendConfig
}

func NewDummyBackend(config *BackendConfig) *DummyBackend {
	backend := &DummyBackend{
		config: config,
	}

	return backend
}

// Listen and publish messages to the `messages` channel so that the app can handle it
func (backend *DummyBackend) Listen(messages chan BackendInputMessage) {
	for {
		sleepTime := rand.Intn(500)
		time.Sleep(time.Duration(sleepTime) * time.Nanosecond)
		messages <- BackendInputMessage{Source: "dummy",
			Content: fmt.Sprintf("{\"test\": true, \"sleep\": %d}", sleepTime)}
	}
}

// Handle a message
func (backend *DummyBackend) HandleMessage(message BackendOutputMessage) error {
	backend.config.Logger.Infow("Handling message", "message", message)
	return nil
}
