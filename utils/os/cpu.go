package osutils

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type CPUUtils interface {
	Usage(ctx context.Context) (float64, error)
}

func NewCPUUitls(ctx context.Context, interval time.Duration) CPUUtils {
	res := &cpuUitls{
		interval: interval,
	}
	_, _ = res.Usage(ctx) // need one more calls when interval == 0
	return res
}

type cpuUitls struct {
	interval time.Duration // sleep interval and cal Usage
}

func (c *cpuUitls) Usage(ctx context.Context) (float64, error) {
	res, err := cpu.PercentWithContext(ctx, c.interval, false) // percpu = false
	if err != nil {
		return 0.0, err
	}
	return res[0], nil
}
