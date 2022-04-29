package search_time

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetYearAndWeek(t *testing.T) {
	testTime := "2006-01-02T15:04:05Z"
	year, week, _ := GetYearAndWeek(testTime)
	assert.Equal(t, year, 2006)
	assert.Equal(t, week, 1)
}

func TestGetYearAndWeekError(t *testing.T) {
	invalidTime := "invalidTime"
	_, _, err := GetYearAndWeek(invalidTime)
	assert.Error(t, err)
}
