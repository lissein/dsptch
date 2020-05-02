package destinations

import "github.com/lissein/dsptch/shared"

type DummyDestination struct {
	config *DestinationConfig
}

func NewDummyDestination(config *DestinationConfig) (*DummyDestination, error) {
	dest := &DummyDestination{config}
	return dest, nil
}

func (dest *DummyDestination) Send(targetIds []int, message shared.Message) error {
	dest.config.logger.Infow("Send message to dummy dest", "message", message)
	return nil
}
