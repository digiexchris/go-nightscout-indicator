package appindicator

import (
	"fmt"
	"github.com/digiexchris/go-nightscout-indicator/icon"
	"github.com/digiexchris/go-nightscout-indicator/unitconverter"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"time"
)

type tray struct {
	units          bool
	errorState     bool
	currentReading unitconverter.Reading //for use if we have to modify the existing display
	newReading     chan unitconverter.Reading
}

type Tray interface {
	Run()
	SetUnits(bool)
	SetTitleValues(reading unitconverter.Reading)
	SetErrorIcon()
	SetNormalIcon()
	GetUnits() bool
	ManageErrorState()
}

func New(r chan unitconverter.Reading) Tray {
	tray := &tray{
		newReading: r,
	}
	tray.SetNormalIcon()
	return tray
}

func (t *tray) GetUnits() bool {
	return t.units
}

func (t *tray) Run() {

	go func() {
		for reading := range t.newReading {
			if reading.Error != nil {
				t.errorState = true
			} else {
				t.errorState = false
			}
			t.SetTitleValues(reading)
		}
	}()

	errorTicker := time.NewTicker(time.Millisecond * 1200)
	go func() {
		for range errorTicker.C {
			t.ManageErrorState()
		}
	}()

	systray.Run(t.onReady, func() {})
}

func (t *tray) SetUnits(mmol bool) {
	t.units = mmol
}

func (t *tray) SetTitleValues(reading unitconverter.Reading) {
	t.currentReading = reading // save for in case we switch units in the future and need to recalculate the current reading
	systray.SetTitle(unitconverter.FormatTitle(t.units, reading.SGV, reading.Delta))
}

func (t *tray) SetErrorIcon() {
	systray.SetIcon(icon.IconError)
}

func (t *tray) SetNormalIcon() {
	systray.SetIcon(icon.IconNormal)
}

func (t *tray) ManageErrorState() {

	if t.currentReading.Error != nil {
		if t.errorState == false {
			t.errorState = true
			t.SetErrorIcon()
		} else {
			t.errorState = false
			t.SetNormalIcon()
		}
	} else {
		if t.errorState == true {
			t.errorState = false
			t.SetNormalIcon()
		}
	}
}

func (t *tray) onReady() {

	systray.SetTitle(unitconverter.FormatTitle(t.units, t.currentReading.SGV, t.currentReading.Delta))
	systray.SetTooltip("Nightscout Indicator")
	mUrl := systray.AddMenuItem("About", "About the indicator")
	mToggle := systray.AddMenuItem("mmol/l", "Switch to mmol/l")

	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Quit", "Quit")

	go func() {
		//todo make these things return events to the statemanager
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		close(t.newReading)
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	for {
		//todo make these things return events to the statemanager
		select {
		case <-mUrl.ClickedCh:
			open.Run("https://github.com/digiexchris")
		case <-mToggle.ClickedCh:

			if t.GetUnits() == unitconverter.MMOL {
				t.SetUnits(unitconverter.MGDL)
				mToggle.SetTitle("mmol/l")
				systray.SetTitle(unitconverter.FormatTitle(t.units, t.currentReading.SGV, t.currentReading.Delta))
			} else {
				t.SetUnits(unitconverter.MMOL)
				mToggle.SetTitle("mg/dl")
				systray.SetTitle(unitconverter.FormatTitle(t.units, t.currentReading.SGV, t.currentReading.Delta))
			}

		}
	}
}
