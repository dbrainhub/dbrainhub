package osutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDiskUsage(t *testing.T) {
	diskUtils := NewDiskUitls()
	usage, err := diskUtils.Usage(context.Background(), "/")
	assert.Nil(t, err)
	assert.Greater(t, usage, 0.0)
}
