package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lissein/dsptch/apps"
	"github.com/lissein/dsptch/backends"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var builtinBackends = map[string]backends.Constructor{
	"redis":     backends.NewRedisBackend,
	"dummy":     backends.NewDummyBackend,
	"websocket": backends.NewWebSocketBackend,
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Dsptch is the main struct for this app
type Dsptch struct {
	backends map[string]backends.Backend
	apps     map[string][]apps.App

	messages chan backends.Message

	logger *zap.SugaredLogger
}

// NewDsptch creates a new app with the specified config
func NewDsptch() (*Dsptch, error) {
	logConfig := zap.NewDevelopmentConfig()
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := logConfig.Build()

	dsptch := &Dsptch{
		logger:   logger.Sugar(),
		backends: make(map[string]backends.Backend),
		messages: make(chan backends.Message, 0),
		apps:     make(map[string][]apps.App),
	}

	dsptch.loadBackends()
	dsptch.loadApps()

	return dsptch, nil
}

func (dsptch *Dsptch) loadBackends() {
	// TODO Load from config
	backendNames := []string{"dummy", "redis", "websocket"}
	backendConfigs := []map[string]interface{}{nil, map[string]interface{}{
		"channels": []string{"test", "blah"},
	}, nil}

	var loadedBackends []string

	for i, backendName := range backendNames {
		dsptch.registerBackend(backendName, &backends.Config{
			Logger: dsptch.logger.Named(backendName),
			Config: backendConfigs[i],
		})
		loadedBackends = append(loadedBackends, backendName)
	}

	dsptch.logger.Info("Backends: ", strings.Join(backendNames, ", "))
}

func (dsptch *Dsptch) registerBackend(name string, config *backends.Config) {
	backendConstructor, found := builtinBackends[name]

	if found {
		backend, err := backendConstructor(config)
		if err != nil {
			dsptch.logger.Panicf("Failed to load backend '%s': %s", name, err)
		}
		dsptch.backends[name] = backend
		return
	}

	// TODO Load backend .so
	dsptch.logger.Panicf("Invalid backend '%s'", name)
}

func (dsptch *Dsptch) loadApps() {
	loadedApps := make([]string, 0)
	err := filepath.Walk("apps/", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".so" {
			return nil
		}

		app, err := LoadApp(path)
		if err != nil {
			return err
		}

		for _, trigger := range app.Triggers() {
			previous, found := dsptch.apps[trigger]

			if !found {
				previous = make([]apps.App, 0)
			}

			dsptch.apps[trigger] = append(previous, app)
			loadedApps = append(loadedApps, fmt.Sprintf("%s[%s]", app.Name(), strings.Join(app.Triggers(), ", ")))
		}
		return nil
	})
	if err != nil {
		dsptch.logger.Panic(err)
	}

	if len(loadedApps) > 0 {
		dsptch.logger.Info("Apps: ", strings.Join(loadedApps, ", "))
	} else {
		dsptch.logger.Info("Apps: None")
	}
}

func (dsptch *Dsptch) Run() error {
	for _, backend := range dsptch.backends {
		go backend.Listen(dsptch.messages)
	}

	// 5 is the number of "workers"
	for i := 0; i < 5; i++ {
		go dsptch.messageHandler(i)
	}

	for {
		fmt.Scanln()
	}
}

func (dsptch *Dsptch) sendApp(backend string, message backends.Message) {
	dest := dsptch.backends[backend]

	if dest == nil {
		dsptch.logger.Panicf("Invalid backend %s", backend)
	}

	dest.Handle(message)
}

func (dsptch *Dsptch) messageHandler(id int) {
	for {
		message := <-dsptch.messages

		dsptch.logger.Infow(fmt.Sprintf("Worker[%d] handling message", id), "message", message)

		apps := dsptch.apps[message.Source]

		for _, app := range apps {
			app.Execute(message, dsptch.sendApp)
		}
	}
}
