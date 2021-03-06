package ali_mns

import (
	"sync/atomic"
)

type QPSMonitor struct {
	latestIndex  int32
	delaySecond  int32
	totalQueries []int32
}

func (p *QPSMonitor) Pulse() {
	index := int32(now().Second()) % p.delaySecond

	if p.latestIndex != index {
		atomic.StoreInt32(&p.latestIndex, index)
		atomic.StoreInt32(&p.totalQueries[p.latestIndex], 0)
	}

	atomic.AddInt32(&p.totalQueries[index], 1)
}

func (p *QPSMonitor) QPS() int32 {
	var totalCount int32 = 0
	for i, queryCount := range p.totalQueries {
		if int32(i) != p.latestIndex {
			totalCount += queryCount
		}
	}
	return totalCount / (p.delaySecond - 1)
}

func NewQPSMonitor(delaySecond int32) *QPSMonitor {
	if delaySecond < 5 {
		delaySecond = 5
	}
	monitor := QPSMonitor{
		delaySecond:  delaySecond,
		totalQueries: make([]int32, delaySecond),
	}
	return &monitor
}
