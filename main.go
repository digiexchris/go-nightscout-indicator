package main

import (
	"github.com/digiexchris/go-nightscout-indicator/applicationmanager"
	"github.com/digiexchris/go-nightscout-indicator/configuration"
)

func main() {
	err := configuration.Load()
	if err != nil {
		panic(err)
	}

	app := applicationmanager.New()
	app.SetUnits(configuration.App.DefaultMmol)

	app.Run()
}
