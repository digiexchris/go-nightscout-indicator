package appindicator

import (
	"fmt"
	"github.com/digiexchris/go-nightscout-indicator/icon"
	"github.com/digiexchris/go-nightscout-indicator/nightscoutclient"
	"github.com/digiexchris/go-nightscout-indicator/unitconverter"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"time"
)

type tray struct {
	units          bool
	errorIconState bool
	oldIconState   bool
	currentReading nightscoutclient.Reading //for use if we have to modify the existing display
	newReading     chan nightscoutclient.Reading
}

type Tray interface {
	Run()
	SetUnits(bool)
	SetTitleValues(reading nightscoutclient.Reading)
	SetErrorIcon()
	SetNormalIcon()
	GetUnits() bool
	ManageErrorState()
}

func New(r chan nightscoutclient.Reading) Tray {
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
				t.errorIconState = true
			} else {
				t.errorIconState = false
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
	systray.SetTitle(unitconverter.FormatTitle(t.units, t.currentReading.SGV, t.currentReading.Delta, t.currentReading.Direction))
}

func (t *tray) SetTitleValues(reading nightscoutclient.Reading) {
	t.currentReading = reading // save for in case we switch units in the future and need to recalculate the current reading
	systray.SetTitle(unitconverter.FormatTitle(t.units, reading.SGV, reading.Delta, reading.Direction))
}

func (t *tray) SetErrorIcon() {
	systray.SetIcon(icon.IconError)
}

func (t *tray) SetOldIcon() {
	systray.SetIcon(icon.IconClock)
}

func (t *tray) SetNormalIcon() {
	systray.SetIcon(icon.IconNormal)
}

func (t *tray) ManageErrorState() {

	if t.currentReading.Error != nil {
		if t.errorIconState == false {
			t.errorIconState = true
			t.SetErrorIcon()
		} else {
			t.errorIconState = false
			t.SetNormalIcon()
		}
	} else {
		if t.errorIconState == true {
			t.errorIconState = false
			t.SetNormalIcon()
		}

		if t.currentReading.OldReading == true {
			if t.oldIconState == false {
				t.oldIconState = true
				t.SetOldIcon()
			} else {
				t.oldIconState = false
				t.SetNormalIcon()
			}
		}
	}
}

func (t *tray) onReady() {

	systray.SetTitle(unitconverter.FormatTitle(t.units, t.currentReading.SGV, t.currentReading.Delta, t.currentReading.Direction))
	systray.SetTooltip("Nightscout Indicator")
	mUrl := systray.AddMenuItem("About", "About the indicator")

	toggleTitle := unitconverter.GetUnitString(unitconverter.MMOL)
	mToggle := systray.AddMenuItem(toggleTitle, fmt.Sprintf("Switch to %s", toggleTitle))

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit")

	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		close(t.newReading)
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	for {
		select {
		case <-mUrl.ClickedCh:
			open.Run("https://github.com/digiexchris")
		case <-mToggle.ClickedCh:

			if t.GetUnits() == unitconverter.MMOL {
				t.SetUnits(unitconverter.MGDL)
				mToggle.SetTitle(unitconverter.GetUnitString(unitconverter.MMOL))
			} else {
				t.SetUnits(unitconverter.MMOL)
				mToggle.SetTitle(unitconverter.GetUnitString(unitconverter.MGDL))
			}

		}
	}
}
