package unitconverter

import "fmt"

const MMOL = true
const MGDL = false

/**
FormatTitle expects values in mg/dl and will output a string appropriate
for the appindicator title
*/
func FormatTitle(units bool, value float32, delta float32) string {

	switch units {
	case MMOL:
		mmolValue := value / 18
		deltaValue := delta / 18
		return fmt.Sprintf("%.1f (%.3f)", mmolValue, deltaValue)
	}

	return fmt.Sprintf("%.0f (%.0f)", value, delta)
}

func GetUnitString(units bool) string {
	if units {
		return "mmol/l"
	} else {
		return "mg/dl"
	}
}
