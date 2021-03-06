package builtins

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lissein/dsptch/backends"
)

type DummyBackend struct {
	config       *backends.Config
	enableListen bool
}

func NewDummyBackend(config *backends.Config) (backends.Backend, error) {
	backend := &DummyBackend{
		config: config,
	}

	return backend, nil
}

// Listen and publish messages to the `messages` channel so that the app can handle it
func (backend *DummyBackend) Listen(messages chan backends.Message) {
	if !backend.enableListen {
		return
	}
	for {
		sleepTime := rand.Intn(500)
		time.Sleep(time.Duration(sleepTime) * time.Nanosecond)
		messages <- backends.Message{Source: "dummy",
			Payload: fmt.Sprintf("{\"test\": true, \"sleep\": %d}", sleepTime)}
	}
}

// Handle a message
func (backend *DummyBackend) Handle(message backends.Message) error {
	backend.config.Logger.Infow("Handling message", "message", message)
	return nil
}
