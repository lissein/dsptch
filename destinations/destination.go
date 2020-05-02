package destinations

import (
	"github.com/lissein/dsptch/shared"
	"go.uber.org/zap"
)

// Destination is the interface for destinations (websockets, push notifications, mails)
type Destination interface {
	Send(targetIds []int, message shared.Message) error
}

type DestinationConfig struct {
	logger *zap.SugaredLogger
}

func NewDestinationConfig(logger *zap.SugaredLogger) *DestinationConfig {
	return &DestinationConfig{logger}
}
