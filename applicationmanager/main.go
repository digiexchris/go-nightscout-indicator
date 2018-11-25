package applicationmanager

import (
	"github.com/digiexchris/go-nightscout-indicator/appindicator"
	"github.com/digiexchris/go-nightscout-indicator/configuration"
	"github.com/digiexchris/go-nightscout-indicator/nightscoutclient"
	"github.com/digiexchris/go-nightscout-indicator/unitconverter"
	"log"
	"time"
)

type ApplicationManager interface {
	RefreshValues()
	SetUnits(defaultMmol bool)
	Run()
}

type applicationManager struct {
	units            bool
	value            float32
	delta            float32
	retrievalError   error
	displayErrorIcon bool
	Readings         chan nightscoutclient.Reading
	AppIndicator     appindicator.Tray
	Client           nightscoutclient.NightscoutClient
}

func (sm *applicationManager) SetUnits(defaultMmol bool) {
	switch defaultMmol {
	case unitconverter.MGDL:
		log.Println("Setting mg/dl")
	case unitconverter.MMOL:
		log.Println("Setting mmol/l")
	}

	sm.units = defaultMmol
	sm.AppIndicator.SetUnits(defaultMmol)
}

func New() ApplicationManager {

	readingChannel := make(chan nightscoutclient.Reading, 4)
	sm := &applicationManager{
		AppIndicator: appindicator.New(readingChannel),
		Readings:     readingChannel,
		Client:       nightscoutclient.New(),
	}
	sm.RefreshValues()

	return sm
}

func (sm *applicationManager) RefreshValues() {
	reading := sm.Client.Get(configuration.App.NightscoutHost, configuration.App.ApiSecret)

	if reading.Error != nil {
		log.Printf("Error reading new values: %s", reading.Error)
	} else {
		oldThreshold := time.Now().Add(-time.Minute * 6)
		if reading.Date.Before(oldThreshold) {
			reading.OldReading = true
		}
		log.Printf("Requested new values: %f %f %s, %s", reading.SGV, reading.Delta, reading.Direction, reading.Date)
	}

	sm.Readings <- reading
}

func (sm *applicationManager) Run() {

	ticker := time.NewTicker(time.Second * configuration.App.CheckInterval)
	go func() {
		for range ticker.C {
			sm.RefreshValues()
		}
	}()

	sm.AppIndicator.Run()
}
