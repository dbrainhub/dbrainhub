package search_time

import "time"

// refer: https://github.com/elastic/kibana/blob/main/src/plugins/vis_types/timelion/common/lib/calculate_interval.ts#L13
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
	switch true {
	case interval <= 500: // <= 0.5s
		return "100ms"
	case interval <= 5000: // <= 5s
		return "1s"
	case interval <= 7500: // <= 7.5s
		return "5s"
	case interval <= 15000: // <= 15s
		return "10s"
	case interval <= 45000: // <= 45s
		return "30s"
	case interval <= 180000: // <= 3m
		return "1m"
	case interval <= 450000: // <= 9m
		return "5m"
	case interval <= 1200000: // <= 20m
		return "10m"
	case interval <= 2700000: // <= 45m
		return "30m"
	case interval <= 7200000: // <= 2h
		return "1h"
	case interval <= 21600000: // <= 6h
		return "3h"
	case interval <= 86400000: // <= 24h
		return "12h"
	case interval <= 604800000: // <= 1w
		return "24h"
	case interval <= 1814400000: // <= 3w
		return "1w"
	case interval < 3628800000: // <  2y
		return "30d"
	default:
		return "1y"
	}
}
