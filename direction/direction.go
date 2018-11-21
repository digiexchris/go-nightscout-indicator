package direction

var directions = map[string]string{
	"NONE":          "⇼",
	"DoubleUp":      "⇈",
	"SingleUp":      "↑",
	"FortyFiveUp":   "↗",
	"Flat":          "→",
	"FortyFiveDown": "↘",
	"SingleDown":    "↓",
	"DoubleDown":    "⇊",
}

//const "NOT COMPUTABLE" = "-"
//const "RATE OUT OF RANGE" = "⇕"

func GetDirectionForTrend(trend string) string {
	if val, ok := directions[trend]; ok {
		return val
	} else {
		return directions["NONE"]
	}
}
