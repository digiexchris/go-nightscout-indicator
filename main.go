package main

import (
	"github.com/digiexchris/go-nightscout-indicator/configuration"
	"github.com/digiexchris/go-nightscout-indicator/statemanager"
)

func main() {
	err := configuration.Load()
	if err != nil {
		panic(err)
	}

	sm := statemanager.New()
	sm.SetUnits(configuration.App.DefaultMmol)

	sm.Run()
}
