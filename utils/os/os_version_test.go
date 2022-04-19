package osutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOsVersion(t *testing.T) {
	var drawin darwinVersionQuerier
	version, err := drawin.GetOsVersion()
	assert.Nil(t, err)
	assert.Equal(t, version, "10.15.7")
}

func TestGetKernalVersion(t *testing.T) {
	var drawin darwinVersionQuerier
	version, err := drawin.GetKernelVersion()
	assert.Nil(t, err)
	assert.Equal(t, version, "19.6.0")
}
