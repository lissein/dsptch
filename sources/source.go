package sources

import (
	"github.com/lissein/dsptch/shared"
	"go.uber.org/zap"
)

// Source is the interface for a source service (Redis, SQS, ...)
type Source interface {
	Listen(messages chan shared.Message) error
}

type SourceConfig struct {
	logger *zap.SugaredLogger
}

func NewSourceConfig(logger *zap.SugaredLogger) *SourceConfig {
	return &SourceConfig{logger}
}
