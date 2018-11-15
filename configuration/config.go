package configuration

import (
	"fmt"
	"github.com/tkanos/gonfig"
	"log"
	"os"
	"path/filepath"
)

var App Config

type Config struct {
	NightscoutHost string
	ApiSecret string
	DefaultMmol bool
}

func Load() error {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	file := fmt.Sprintf("%s/config.json",dir)
	log.Printf("Attempting to load configuration from %s", file)

	App = Config{}
	return gonfig.GetConf(file, &App)
}

