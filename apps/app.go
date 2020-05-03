package apps

import "github.com/lissein/dsptch/backends"

type InitAppFunction = func() (App, error)

type SendFunction = func(backend string, message backends.Message)

type App interface {
	Name() string
	Triggers() []string
	Execute(message backends.Message, send SendFunction)
}
