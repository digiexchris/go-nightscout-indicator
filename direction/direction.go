package direction

var directions = map[string]string{
	"NONE":              "⇼",
	"DoubleUp":          "⇈",
	"SingleUp":          "↑",
	"FortyFiveUp":       "↗",
	"Flat":              "→",
	"FortyFiveDown":     "↘",
	"SingleDown":        "↓",
	"DoubleDown":        "⇊",
	"NOT COMPUTABLE":    "-",
	"RATE OUT OF RANGE": "⇕",
}

//todo make this modifyable by config file so people can put custom emoji in here or something fun

func GetDirectionForTrend(trend string) string {
	if val, ok := directions[trend]; ok {
		return val
	} else {
		return directions["NONE"]
	}
}
