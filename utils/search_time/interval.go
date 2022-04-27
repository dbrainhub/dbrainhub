package search_time

import "time"

var IntervalSeconds = []float64{
	500,
	5000,
	7500,
	15000,
	45000,
	180000,
	450000,
	1200000,
	2700000,
	7200000,
	21600000,
	86400000,
	604800000,
	1814400000,
	3628800000,
}
var IntervalStrings = []string{
	"100ms",
	"1s",
	"5s",
	"10s",
	"30s",
	"1m",
	"5m",
	"10m",
	"30m",
	"1h",
	"3h",
	"12h",
	"24h",
	"1w",
	"30d",
}

func GetInterval(from, to string) string {
	fromT, err := time.ParseInLocation("2006-01-02T15:04:05Z", from, time.Local)
	if err != nil {
		return "1y"
	}
	toT, err := time.ParseInLocation("2006-01-02T15:04:05Z", to, time.Local)
	if err != nil {
		return "1y"
	}
	interval := toT.Sub(fromT).Seconds()

	intervalStr := "1y"
	loc := -1
	for i, v := range IntervalSeconds {
		if interval <= v {
			loc = i
			break
		}
	}
	if loc > 0 {
		intervalStr = IntervalStrings[loc]
	}

	return intervalStr
}
