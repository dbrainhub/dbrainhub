package search_time

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterval(t *testing.T) {
	from := "2021-04-06T11:43:47.613Z"
	to := "2021-04-25T07:48:47.613Z"
	interval := GetInterval(from, to)
	assert.Equal(t, interval, "30m")
}

func TestIntervalDefault(t *testing.T) {
	from := ""
	to := ""
	interval := GetInterval(from, to)
	assert.Equal(t, interval, "1y")
}
