package main

import (
	"errors"
	"fmt"
	"github.com/digiexchris/go-nightscout-indicator/configuration"
	"github.com/digiexchris/go-nightscout-indicator/icon"
	"github.com/skratchdot/open-golang/open"
	"log"
	"time"

	"github.com/getlantern/systray"
)

const MMOL = true
const MGDL = false

var value float32
var delta float32
var units bool
var retrievalError error
var errorIcon bool

func main() {
	err := configuration.Load()
	if err != nil {
		panic(err)
	}

	systray.Run(onReady, func() {})
}

func onReady() {

	value = 100
	delta = -4

	errorIcon = false
	retrievalError = nil

	systray.SetIcon(icon.IconNormal)

	value, delta, retrievalError = refresh()

	units = configuration.App.DefaultMmol

	switch (units) {
	case MGDL:
		log.Println("Setting mg/dl")
	case MMOL:
		log.Println("Setting mmol/l")
	default:
		panic(errors.New("Invalid value for DefaultMmol in config.json"))
	}

	systray.SetTitle(format(units, value, delta))
	systray.SetTooltip("Nightscout Indicator")
	mUrl := systray.AddMenuItem("About", "About the indicator")
	mToggle := systray.AddMenuItem("mmol/l", "Switch to mmol/l")

	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Quit", "Quit")

	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	errorTicker := time.NewTicker(time.Millisecond * 1200)

	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for range ticker.C {
			value, delta, retrievalError = refresh()

			if retrievalError != nil {
				//show a red icon
				log.Println("Error retrieving data")
			} else {
				//show a green icon
				retrievalError = nil
				systray.SetTitle(format(units, value, delta))
			}
		}
	}()

	go func() {
		for range errorTicker.C {
			if retrievalError != nil {
				if errorIcon == false {
					errorIcon = true
					systray.SetIcon(icon.IconError)
				} else {
					errorIcon = false
					systray.SetIcon(icon.IconNormal)
				}
			}
		}
	}()

	for {
		select {
		case <-mUrl.ClickedCh:
			open.Run("https://github.com/digiexchris")
		case <-mToggle.ClickedCh:
			if units == MMOL {
				units = MGDL
				mToggle.SetTitle("mmol/l")
				systray.SetTitle(format(units, value, delta))
			} else {
				units = MMOL
				mToggle.SetTitle("mg/dl")
				systray.SetTitle(format(units, value, delta))
			}

		}
	}
}

func format(units bool, value float32, delta float32, ) string {

	switch units {
	case MMOL:
		mmolValue := value / 18
		deltaValue := delta / 18
		return fmt.Sprintf("%.1f (%.3f) ", mmolValue, deltaValue)
	}

	return fmt.Sprintf("%.0f (%.0f) ", value, delta)
}
