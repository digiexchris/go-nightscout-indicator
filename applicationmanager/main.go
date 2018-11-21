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
	Readings         chan unitconverter.Reading
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

	readingChannel := make(chan unitconverter.Reading, 4)
	sm := &applicationManager{
		AppIndicator: appindicator.New(readingChannel),
		Readings:     readingChannel,
		Client:       nightscoutclient.New(),
	}
	sm.Readings <- sm.Client.Get(configuration.App.NightscoutHost, configuration.App.ApiSecret)

	return sm
}

func (sm *applicationManager) RefreshValues() {
	reading := sm.Client.Get(configuration.App.NightscoutHost, configuration.App.ApiSecret)
	sm.Readings <- reading
	if reading.Error != nil {
		log.Printf("Error reading new values: %s", reading.Error)
	} else {
		log.Printf("Requested new values: %f %f", reading.SGV, reading.Delta)
	}
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
