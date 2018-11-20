package applicationmanager

import (
	"encoding/json"
	"fmt"
	"github.com/digiexchris/go-nightscout-indicator/appindicator"
	"github.com/digiexchris/go-nightscout-indicator/configuration"
	"github.com/digiexchris/go-nightscout-indicator/unitconverter"
	"io/ioutil"
	"log"
	"net/http"
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
	}
	sm.Readings <- refresh()

	return sm
}

func (sm *applicationManager) RefreshValues() {
	reading := refresh()
	sm.Readings <- reading
	if reading.Error != nil {
		log.Printf("Error reading new values: %s", reading.Error)
	} else {
		log.Printf("Requested new values: %f %f", reading.SGV, reading.Delta)
	}
}

func refresh() (reading unitconverter.Reading) {

	type values struct {
		Sgv   float32
		Delta float32
	}

	url := fmt.Sprintf("https://%s/api/v1/entries/current.json", configuration.App.NightscoutHost)
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("API-SECRET", configuration.App.ApiSecret)

	response, err := client.Do(req)
	if err != nil {
		return unitconverter.Reading{
			Error: err,
		}
	} else {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return unitconverter.Reading{
				Error: err,
			}
		}

		var v []values
		err = json.Unmarshal(data, &v)
		if err != nil {
			return unitconverter.Reading{
				Error: err,
			}
		}

		return unitconverter.Reading{
			SGV:   v[0].Sgv,
			Delta: v[0].Delta,
			Error: nil,
		}
	}

	return unitconverter.Reading{
		Error: err,
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
