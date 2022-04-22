package osutils

import (
	"context"

	"github.com/shirou/gopsutil/mem"
)

type MemoryUtils interface {
	Usage(ctx context.Context) (float64, error)
}

func NewMemoryUitls() MemoryUtils {
	return &memoryUitls{}
}

type memoryUitls struct {
}

func (m *memoryUitls) Usage(ctx context.Context) (float64, error) {
	memInfo, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return 0.0, err
	}
	return memInfo.UsedPercent, nil
}
