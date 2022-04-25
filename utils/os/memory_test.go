package osutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMemory(t *testing.T) {
	memoryUtils := NewMemoryUitls()
	usage, err := memoryUtils.Usage(context.Background())
	assert.Nil(t, err)
	assert.Greater(t, usage, 0.0)
}
