package main

import (
	"github.com/lissein/dsptch/internal"
)

func main() {
	dsptch, err := internal.NewDsptch()

	if err != nil {
		panic(err)
	}

	dsptch.Run()
}
