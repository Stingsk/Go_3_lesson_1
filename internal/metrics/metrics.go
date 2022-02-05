package metrics

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

var jsonMonitor []byte

type Monitor struct {
	Alloc,
	BuckHashSys,
	Frees,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	OtherSys,
	PauseTotalNs,
	StackInuse,
	StackSys,
	RandomValue,
	Sys uint64
	NumGC,
	NumForcedGC uint32
	GCCPUFraction float64
	NumGoroutine,
	PollCount int
}

func NewMonitor(duration int) {
	var m Monitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	count := 0
	for {
		<-time.After(interval)

		count++
		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.BuckHashSys = rtm.BuckHashSys
		m.Frees = rtm.Frees
		m.GCCPUFraction = rtm.GCCPUFraction
		m.GCSys = rtm.GCSys
		m.HeapAlloc = rtm.HeapAlloc
		m.HeapIdle = rtm.HeapIdle
		m.HeapInuse = rtm.HeapInuse
		m.HeapObjects = rtm.HeapObjects
		m.HeapReleased = rtm.HeapReleased
		m.HeapSys = rtm.HeapSys
		m.LastGC = rtm.LastGC
		m.Lookups = rtm.Lookups
		m.MCacheInuse = rtm.MCacheInuse
		m.MCacheSys = rtm.MCacheSys
		m.MSpanInuse = rtm.MSpanInuse
		m.MSpanSys = rtm.MSpanSys
		m.Mallocs = rtm.Mallocs
		m.NextGC = rtm.NextGC
		m.NumForcedGC = rtm.NumForcedGC
		m.NumGC = rtm.NumGC
		m.OtherSys = rtm.OtherSys
		m.PauseTotalNs = rtm.PauseTotalNs
		m.StackInuse = rtm.StackInuse
		m.StackSys = rtm.StackSys
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		m.PollCount = count
		m.RandomValue = rand.Uint64()

		// Live objects = Mallocs - Frees

		// Just encode to json and print
		jsonMonitor, _ = json.Marshal(m)
	}
}

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	w.Write(jsonMonitor)
}
