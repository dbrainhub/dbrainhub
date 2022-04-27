package search_time

import "time"

var Intervals = []Interval{
	{500, "100ms"},
	{5000, "1s"},
	{7500, "5s"},
	{15000, "10s"},
	{45000, "30s"},
	{180000, "1m"},
	{450000, "5m"},
	{1200000, "10m"},
	{2700000, "30m"},
	{7200000, "1h"},
	{21600000, "3h"},
	{86400000, "12h"},
	{604800000, "24h"},
	{1814400000, "1w"},
	{3628800000, "30d"},
}

type Interval struct {
	IntervalMs     int64
	IntervalString string
}

func GetInterval(from, to string, buckets int64) string {
	intervalStr := "1y"
	fromT, err := time.ParseInLocation("2006-01-02T15:04:05Z", from, time.Local)
	if err != nil {
		return intervalStr
	}
	toT, err := time.ParseInLocation("2006-01-02T15:04:05Z", to, time.Local)
	if err != nil {
		return intervalStr
	}
	interval := toT.Sub(fromT).Milliseconds() / buckets

	for _, v := range Intervals {
		if interval <= v.IntervalMs {
			intervalStr = v.IntervalString
			break
		}
	}

	return intervalStr
}
