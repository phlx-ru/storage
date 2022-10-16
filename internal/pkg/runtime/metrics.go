package runtime

import (
	"context"
	"runtime"
	"runtime/debug"
	"time"

	"storage/internal/pkg/metrics"
)

const (
	maxGCPausesTimings = 10
)

func CollectGoMetrics(ctx context.Context, metric metrics.Metrics) {
	collect := goCollector(metric)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		collect()
		time.Sleep(time.Second)
	}
}

func goCollector(metrics metrics.Metrics) func() {
	lastNumGC := int64(0)

	var statGC debug.GCStats
	var statMem runtime.MemStats

	return func() {

		// goroutines and threads
		metrics.Gauge(`go.goroutines`, runtime.NumGoroutine())
		n, _ := runtime.ThreadCreateProfile(nil)
		metrics.Gauge(`go.threads`, n)

		// garbage collector stats
		debug.ReadGCStats(&statGC)
		pausesCount := int(statGC.NumGC - lastNumGC)
		if pausesCount > len(statGC.Pause) {
			pausesCount = len(statGC.Pause) // 256*2+3
		}
		for i := 0; i < pausesCount && i < maxGCPausesTimings; i++ {
			metrics.Timing(
				`go.gc_pause_microseconds`,
				float64(statGC.Pause[i]*time.Microsecond)/float64(time.Millisecond),
			)
		}
		lastNumGC = statGC.NumGC

		// memory stats
		runtime.ReadMemStats(&statMem)
		metrics.Gauge(`go.mem_alloc_bytes`, statMem.Alloc)
		metrics.Gauge(`go.mem_alloc_bytes_total`, statMem.TotalAlloc)
		metrics.Gauge(`go.mem_sys_bytes`, statMem.Sys)
		metrics.Gauge(`go.mem_heap_alloc_bytes`, statMem.HeapAlloc)
	}
}
