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

func (source *DummySource) Listen(messages chan shared.SourceMessage) error {
	for {
		sleepTime := rand.Intn(500)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		messages <- shared.SourceMessage{Source: "dummy_src",
			Content: map[string]interface{}{"test": true, "sleep": sleepTime}}
	}
}
