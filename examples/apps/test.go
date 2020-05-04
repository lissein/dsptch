package main

import (
	"github.com/lissein/dsptch/apps"
	"github.com/lissein/dsptch/backends"
	"github.com/lissein/dsptch/builtins"
)

type TestApp struct {
}

func InitApp() (apps.App, error) {
	return &TestApp{}, nil
}

func (app *TestApp) Name() string {
	return "test-app"
}

func (app *TestApp) Triggers() []string {
	return []string{"websocket"}
}

func (app *TestApp) Execute(message backends.Message, send apps.SendFunction) {
	targets := make([]int, 1)
	targets[0] = 2
	send("websocket", backends.Message{
		Source: message.Source,
		Payload: &builtins.WebSocketHandlePayload{
			Targets: targets,
		},
	})
}
