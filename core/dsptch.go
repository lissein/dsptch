package core

import (
	"fmt"

	"github.com/lissein/dsptch/destinations"
	"github.com/lissein/dsptch/shared"
	"github.com/lissein/dsptch/sources"
	"github.com/lissein/dsptch/storages"
	"go.uber.org/zap"
)

// Dsptch is the main struct for this app
type Dsptch struct {
	destinations []destinations.Destination
	sources      []sources.Source
	storage      storages.Storage

	logger   *zap.SugaredLogger
	messages chan shared.Message
}

// NewDsptch creates a new app with the specified config
func NewDsptch() (*Dsptch, error) {
	logger, _ := zap.NewDevelopment()
	dsptch := &Dsptch{
		logger:       logger.Sugar(),
		destinations: make([]destinations.Destination, 0),
		sources:      make([]sources.Source, 0),
		messages:     make(chan shared.Message, 0),
	}

	// Register destinations
	destConfig := destinations.NewDestinationConfig(dsptch.logger)
	dummyDest, err := destinations.NewDummyDestination(destConfig)
	if err != nil {
		panic(err)
	}

	dsptch.destinations = append(dsptch.destinations, dummyDest)

	// Register sources
	sourceConfig := sources.NewSourceConfig(dsptch.logger)
	dummySource, err := sources.NewDummySource(sourceConfig)
	if err != nil {
		panic(err)
	}

	dsptch.sources = append(dsptch.sources, dummySource)
	return dsptch, nil
}

func (dsptch *Dsptch) Run() error {
	for _, source := range dsptch.sources {
		go source.Listen(dsptch.messages)
	}

	// 5 is the number of "workers"
	for i := 0; i < 5; i++ {
		go dsptch.messageHandler(i)
	}

	for {
		fmt.Scanln()
	}
}

func (dsptch *Dsptch) messageHandler(id int) {
	dsptch.logger.Infof("Started worker %d", id)

	for {
		message := <-dsptch.messages

		dsptch.logger.Infow(fmt.Sprintf("Worker[%d] handling message", id), "message", message)

		// Execute tengo source with context (message, storage, available destinations, ...)
		// And get results (target ids, updated message)

		targetIds := []int{0, 1, 2, 3}

		dest := dsptch.destinations[0]
		dest.Send(targetIds, message)
	}
}
