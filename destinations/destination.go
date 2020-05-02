package destinations

import (
	"go.uber.org/zap"
)

// Destination is the interface for destinations (websockets, push notifications, mails)
type Destination interface {
	Send(message interface{}) error
}

type DestinationConfig struct {
	logger *zap.SugaredLogger
}

func NewDestinationConfig(logger *zap.SugaredLogger) *DestinationConfig {
	return &DestinationConfig{logger}
}
