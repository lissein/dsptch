package main

import (
	"github.com/lissein/dsptch/core"
)

func main() {
	dsptch, err := core.NewDsptch()

	if err != nil {
		panic(err)
	}

	dsptch.Run()
}
