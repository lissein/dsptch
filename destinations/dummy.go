package destinations

type DummyDestination struct {
	config *DestinationConfig
}

func NewDummyDestination(config *DestinationConfig) (*DummyDestination, error) {
	dest := &DummyDestination{config}
	return dest, nil
}

func (dest *DummyDestination) Send(message interface{}) error {
	dest.config.logger.Infow("Send message to dummy dest", "message", message)
	return nil
}
