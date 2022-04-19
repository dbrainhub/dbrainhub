package osutils

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOsVersion(t *testing.T) {
	querier := NewVersionQuerier(runtime.GOOS)
	_, err := querier.GetOsVersion()
	assert.Nil(t, err)
}

func TestGetKernalVersion(t *testing.T) {
	querier := NewVersionQuerier(runtime.GOOS)
	_, err := querier.GetKernelVersion()
	assert.Nil(t, err)
}
