package configuration

import (
	"fmt"
	"github.com/tkanos/gonfig"
	"log"
	"os"
	"time"
)

var App Config

type Config struct {
	NightscoutHost string
	ApiSecret      string
	DefaultMmol    bool
	CheckInterval  time.Duration
}

func Load() error {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file := fmt.Sprintf("%s/config.json", dir)
	log.Printf("Attempting to load configuration from %s", file)

	App = Config{}
	err = gonfig.GetConf(file, &App)

	log.Printf("Refresh Interval: %d", App.CheckInterval)

	return err
}
