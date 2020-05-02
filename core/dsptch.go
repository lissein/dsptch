package core

import (
	"fmt"
	"io/ioutil"

	"github.com/d5/tengo/v2"
	"github.com/lissein/dsptch/destinations"
	"github.com/lissein/dsptch/shared"
	"github.com/lissein/dsptch/sources"
	"github.com/lissein/dsptch/storages"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Dsptch is the main struct for this app
type Dsptch struct {
	destinations map[string]destinations.Destination
	sources      map[string]sources.Source
	storage      storages.Storage
	scripts      map[string]*tengo.Compiled

	logger   *zap.SugaredLogger
	messages chan shared.SourceMessage
}

// NewDsptch creates a new app with the specified config
func NewDsptch() (*Dsptch, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()

	dsptch := &Dsptch{
		logger:       logger.Sugar(),
		destinations: make(map[string]destinations.Destination),
		sources:      make(map[string]sources.Source),
		messages:     make(chan shared.SourceMessage, 0),
		scripts:      make(map[string]*tengo.Compiled),
	}

	// Register destinations
	destConfig := destinations.NewDestinationConfig(dsptch.logger)
	dummyDest, err := destinations.NewDummyDestination(destConfig)
	if err != nil {
		panic(err)
	}
	dsptch.destinations["dummy"] = dummyDest

	// Register sources
	sourceConfig := sources.NewSourceConfig(dsptch.logger)
	dummySource, err := sources.NewDummySource(sourceConfig)
	if err != nil {
		panic(err)
	}

	dsptch.sources["dummy"] = dummySource

	// Load scripts
	content, err := ioutil.ReadFile("scripts/dummy_src.tengo")
	if err != nil {
		panic(err)
	}

	script := tengo.NewScript(content)
	script.Add("destinations", nil)
	script.Add("storage", nil)
	script.Add("input", nil)

	compiled, err := script.Compile()
	if err != nil {
		panic(err)
	}
	dsptch.scripts["dummy_src"] = compiled

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

		// targetIds := []int{0, 1, 2, 3}

		script := dsptch.scripts[message.Source]

		if script != nil {
			err := script.Set("input", message.Content)
			if err != nil {
				panic(err)
			}
		}

		if err := script.Run(); err != nil {
			panic(err)
		}

		targetIds := toIntSlice(script.Get("targets").Array())
		destID := script.Get("destination").String()
		destMessage := shared.DestinationMessage{
			Source:  message.Source,
			Content: script.Get("output").Map(),
		}

		dest := dsptch.destinations[destID]

		if dest == nil {
			dsptch.logger.Panicw("Invalid destination", "destination", destID)
		}

		dest.Send(targetIds, destMessage)
	}
}
