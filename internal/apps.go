package internal

import (
	"errors"
	"fmt"
	"plugin"

	"github.com/lissein/dsptch/apps"
)

func LoadApp(path string) (apps.App, error) {
	plug, err := plugin.Open(path)

	if err != nil {
		return nil, err
	}

	initAppSym, err := plug.Lookup("InitApp")
	if err != nil {
		return nil, err
	}

	initApp, ok := initAppSym.(apps.InitAppFunction)
	if !ok {
		return nil, errors.New("App doesn't have a valid 'InitApp' function")
	}

	loadedApp, err := initApp()
	if err != nil {
		return nil, fmt.Errorf("Failed to load app: %s", err)
	}

	return loadedApp, nil
}
