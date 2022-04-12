package rate_limit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlidingWindow(t *testing.T) {
	// 时间窗口覆盖 30ms
	sw := NewSlidingWindow(3, 10)
	now := int64(0)
	const metric = "success"

	for i := 0; i < 40; i++ {
		sw.Hit(now, metric)
		now++
	}
	// 40ms 内，每 ms Hit 一次，但滑动窗口只会统计最近 30ms 的数据
	assert.Equal(t, sw.GetHit(now, metric), int64(30))

	// 这里体现滑动窗口的统计偏差(cell 时间范围越小，cell 数量越多，才越精确)。
	// 第 42ms 时，滑动窗口里各的 cell 的时间范围是：
	// [20, 30): 10 次
	// [30, 40): 10 次
	// [40, 50):  0 次
	now = 42
	assert.Equal(t, sw.GetHit(now, metric), int64(20))
}
