package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"
)

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

func RunGetMetrics(ctx context.Context, duration int, messages *SensorData, wg *sync.WaitGroup) error {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(duration) * time.Second) // создаём таймер
	count := 0
	for {
		count++
		select {
		case <-ticker.C:
			metrics := GetMetrics(count)
			messages.Store(metrics)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func GetMetrics(count int) []string {
	monitor, err := metrics.NewMonitor(count)
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
		switch value.Kind() {
		case reflect.Uint64, reflect.Uint32:
			val = strconv.FormatUint(value.Uint(), 10)
		case reflect.Int:
			val = strconv.FormatInt(value.Int(), 10)
		case reflect.Float64:
			val = strconv.FormatFloat(value.Float(), 'f', 6, 64)
		}
		nameType := value.Type().Name()
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

func main() {
	var wg sync.WaitGroup
	var sensorData SensorData
	sensorData.Store(GetNames())
	ctx, _ := context.WithCancel(context.Background())

	wg.Add(1)
	go RunGetMetrics(ctx, 2, &sensorData, &wg)

	wg.Add(2)
	go RunSender(ctx, 10, &sensorData, &wg)

	wg.Wait()
}

func RunSender(ctx context.Context, duration int, messages *SensorData, wg *sync.WaitGroup) error {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(duration) * time.Second) // создаём таймер
	for {
		select {
		case <-ticker.C:
			messagesFromChan := messages.Get()
			for _, mes := range messagesFromChan {
				Send(mes)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
func Send(send string) {
	endpoint := "http://localhost:8080/update/" + send
	// конструируем HTTP-клиент
	client := &http.Client{}
	// конструируем запрос
	// запрос методом POST должен, кроме заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	// в большинстве случаев отлично подходит bytes.Buffer
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString("monitor"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Content-Length", strconv.Itoa(len("monitor")))
	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// печатаем код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// и печатаем его
	fmt.Println(string(body))
}
