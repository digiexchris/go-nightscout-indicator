package unitconverter

import (
	"fmt"
	direction2 "github.com/digiexchris/go-nightscout-indicator/direction"
)

const MMOL = true
const MGDL = false

/**
FormatTitle expects values in mg/dl and will output a string appropriate
for the appindicator title
*/
func FormatTitle(units bool, value float32, delta float32, direction string) string {

	trend := direction2.GetDirectionForTrend(direction)
	switch units {
	case MMOL:
		mmolValue := value / 18
		deltaValue := delta / 18
		return fmt.Sprintf("%.1f (%.3f %s)", mmolValue, deltaValue, trend)
	}

	return fmt.Sprintf("%.0f (%.0f %s)", value, delta, trend)
}

func GetUnitString(units bool) string {
	if units {
		return "mmol/l"
	} else {
		return "mg/dl"
	}
}
