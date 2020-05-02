package sources

import (
	"math/rand"
	"time"

	"github.com/lissein/dsptch/shared"
)

type DummySource struct {
	config *SourceConfig
}

func NewDummySource(config *SourceConfig) (*DummySource, error) {
	source := &DummySource{config}
	return source, nil
}

func (source *DummySource) Listen(messages chan shared.Message) error {
	for {
		sleepTime := rand.Intn(500)
		time.Sleep(time.Duration(sleepTime) * time.Microsecond)
		messages <- shared.Message{"test": true, "sleep": sleepTime}
	}
}
