package core

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/lissein/dsptch/backends"
	"github.com/lissein/dsptch/storages"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Map backend name to backend constructor
var backendConstructors = map[string]interface{}{
	"redis": backends.NewRedisBackend,
	"dummy": backends.NewDummyBackend,
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Dsptch is the main struct for this app
type Dsptch struct {
	scripts  map[string]*tengo.Compiled
	backends map[string]backends.Backend
	messages chan backends.BackendInputMessage

	storage storages.Storage

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
		messages: make(chan backends.BackendInputMessage, 0),
		scripts:  make(map[string]*tengo.Compiled),
	}

	dsptch.loadBackends()
	dsptch.loadScripts()

	return dsptch, nil
}

func (dsptch *Dsptch) loadBackends() {
	// TODO Load from config
	backendNames := []string{"dummy"}
	var loadedBackends []string

	for _, backendName := range backendNames {
		dsptch.registerBackend(backendName)
		loadedBackends = append(loadedBackends, backendName)
	}

	dsptch.logger.Info("Backends: ", strings.Join(backendNames, ", "))
}

func (dsptch *Dsptch) registerBackend(name string) {
	if name == "dummy" {
		dsptch.backends[name] = backends.NewDummyBackend(&backends.BackendConfig{
			Logger: dsptch.logger,
		})
		return
	}

	dsptch.logger.Panicf("Invalid backend '%s'", name)
}

func (dsptch *Dsptch) loadScripts() {
	// redisScript := dsptch.loadScript("scripts/test.tengo")

	// dsptch.scripts["redis/test"] = redisScript
	// dsptch.scripts["redis/blah"] = redisScript

	dummyScript := dsptch.loadScript("scripts/dummy.tengo")
	dsptch.scripts["dummy"] = dummyScript
}

func (dsptch *Dsptch) loadScript(filename string) *tengo.Compiled {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		dsptch.logger.Panic(err)
	}

	script := tengo.NewScript(content)

	script.SetImports(stdlib.GetModuleMap("json"))
	Must(script.Add("storage", nil))
	Must(script.Add("input", nil))
	Must(script.Add("send", &tengo.UserFunction{
		Name:  "send",
		Value: dsptch.send,
	}))

	compiled, err := script.Compile()
	if err != nil {
		dsptch.logger.Panic(err)
	}
	return compiled
}

func (dsptch *Dsptch) Run() error {
	for _, backend := range dsptch.backends {
		go backend.Listen(dsptch.messages)
	}

	// 5 is the number of "workers"
	for i := 0; i < 5; i++ {
		clonedScripts := make(map[string]*tengo.Compiled)

		for k, v := range dsptch.scripts {
			clonedScripts[k] = v.Clone()
		}

		go dsptch.messageHandler(i, clonedScripts)
	}

	for {
		fmt.Scanln()
	}
}

func (dsptch *Dsptch) send(args ...tengo.Object) (tengo.Object, error) {
	destID := (args[0].(*tengo.String)).Value
	message := (args[1].(*tengo.String)).Value

	dest := dsptch.backends[destID]

	if dest == nil {
		dsptch.logger.Panicf("Invalid destination %s", destID)
	}

	dest.HandleMessage(backends.BackendOutputMessage{Content: message})
	return nil, nil
}

func (dsptch *Dsptch) messageHandler(id int, scripts map[string]*tengo.Compiled) {
	dsptch.logger.Infof("Started worker %d", id)

	for {
		message := <-dsptch.messages

		dsptch.logger.Infow(fmt.Sprintf("Worker[%d] handling message", id), "message", message)

		// TODO handle list of scenario per sources
		script := scripts[message.Source]

		if script == nil {
			dsptch.logger.Panicf("Script %s not found", message.Source)
		}

		err := script.Set("input", message.Content)
		if err != nil {
			dsptch.logger.Panic(err)
		}

		if err := script.Run(); err != nil {
			dsptch.logger.Panic(err)
		}
	}
}
