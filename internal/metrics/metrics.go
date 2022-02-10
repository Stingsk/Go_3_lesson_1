package metrics

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Monitor struct {
	Alloc         uint64  `json:"Alloc"`
	BuckHashSys   uint64  `json:"BuckHashSys"`
	Frees         uint64  `json:"Frees"`
	GCSys         uint64  `json:"GCSys"`
	HeapAlloc     uint64  `json:"HeapAlloc"`
	HeapIdle      uint64  `json:"HeapIdle"`
	HeapInuse     uint64  `json:"HeapInuse"`
	HeapObjects   uint64  `json:"HeapObjects"`
	HeapReleased  uint64  `json:"HeapReleased"`
	HeapSys       uint64  `json:"HeapSys"`
	LastGC        uint64  `json:"LastGC"`
	Lookups       uint64  `json:"Lookups"`
	MCacheInuse   uint64  `json:"MCacheInuse"`
	MCacheSys     uint64  `json:"MCacheSys"`
	MSpanInuse    uint64  `json:"MSpanInuse"`
	MSpanSys      uint64  `json:"MSpanSys"`
	Mallocs       uint64  `json:"Mallocs"`
	NextGC        uint64  `json:"NextGC"`
	OtherSys      uint64  `json:"OtherSys"`
	PauseTotalNs  uint64  `json:"PauseTotalNs"`
	StackInuse    uint64  `json:"StackInuse"`
	StackSys      uint64  `json:"StackSys"`
	RandomValue   uint64  `json:"RandomValue"`
	Sys           uint64  `json:"Sys"`
	NumGC         uint32  `json:"NumGC"`
	NumForcedGC   uint32  `json:"NumForcedGC"`
	GCCPUFraction float64 `json:"GCCPUFraction"`
	NumGoroutine  int     `json:"NumGoroutine"`
	PollCount     int     `json:"PollCount"`
}

type SensorData struct {
	mu   sync.RWMutex
	last []string
}

func (d *SensorData) Store(data []string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.last = data
}

func (d *SensorData) Get() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.last
}

func RunGetMetrics(ctx context.Context, duration int, messages *SensorData, wg *sync.WaitGroup, sigChan chan os.Signal) error {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(duration) * time.Second) // создаём таймер
	count := 0
	for {
		count++
		select {
		case <-ticker.C:
			metrics := getMetrics(count)
			messages.Store(metrics)
		case <-ctx.Done():
			return ctx.Err()
		case <-sigChan:
			return errors.New("аварийное завершение")
		}
	}
}

func getMetrics(count int) []string {
	monitor, err := newMonitor(count)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	v := reflect.ValueOf(monitor)

	names := GetNames()

	result := make([]string, len(names))
	for i, name := range names {
		value := v.FieldByName(name)
		val := ""
		nameType := ""
		switch value.Kind() {
		case reflect.Uint64, reflect.Uint32:
			val = strconv.FormatUint(value.Uint(), 10)
			nameType = "gauge"
		case reflect.Int:
			val = strconv.FormatInt(value.Int(), 10)
			nameType = "counter"
		case reflect.Float64:
			val = strconv.FormatFloat(value.Float(), 'f', 6, 64)
			nameType = "gauge"
		}
		result[i] = nameType + "/" + name + "/" + val
	}

	return result
}

func GetNames() []string {
	result := []string{"Alloc",
		"BuckHashSys",
		"Frees",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"RandomValue",
		"Sys",
		"NumGC",
		"NumForcedGC",
		"GCCPUFraction",
		"NumGoroutine",
		"PollCount"}

	return result
}

func newMonitor(count int) (Monitor, error) {
	var m Monitor
	var rtm runtime.MemStats
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

	return m, nil
}
