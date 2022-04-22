package osutils

import (
	"context"

	"github.com/shirou/gopsutil/disk"
)

type DiskUtils interface {
	Usage(ctx context.Context, path string) (float64, error)
}

func NewDiskUitls() DiskUtils {
	return &diskUitls{}
}

type diskUitls struct {
}

func (m *diskUitls) Usage(ctx context.Context, path string) (float64, error) {
	stat, err := disk.UsageWithContext(ctx, path)
	if err != nil {
		return 0.0, err
	}
	return stat.UsedPercent, nil
}
