package sources

import (
	"fmt"
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
		time.Sleep(time.Duration(sleepTime) * time.Nanosecond)
		messages <- shared.SourceMessage{Source: "dummy_src",
			Content: fmt.Sprintf("{\"test\": true, \"sleep\": %d}", sleepTime)}
	}
}
