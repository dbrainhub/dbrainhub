package osutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCPU(t *testing.T) {
	cpuUtils := NewCPUUitls(context.Background(), 0)
	usage, err := cpuUtils.Usage(context.Background())
	assert.Nil(t, err)
	assert.Greater(t, usage, 0.0)
}
