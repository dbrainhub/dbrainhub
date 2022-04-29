package search_time

import "time"

func GetYearAndWeek(datetime string) (y, w int, err error) {
	tmp, err := time.ParseInLocation("2006-01-02T15:04:05Z", datetime, time.Local)
	if err != nil {
		return 0, 0, err
	}
	year, week := tmp.ISOWeek()
	return year, week, nil
}
