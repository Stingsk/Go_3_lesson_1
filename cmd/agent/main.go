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
	"strings"
	"time"
)

func RunGetMetrics(ctx context.Context, duration int, messages chan []string) error {
	ticker := time.NewTicker(time.Duration(duration) * time.Second) // создаём таймер
	count := 0
	for {
		count++
		select {
		case <-ticker.C:
			messages <- GetMetrics(count)
			return nil
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

	result := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		val := ""
		json := ""
		switch v.Field(i).Kind() {
		case reflect.Uint64, reflect.Uint32:
			val = strconv.FormatUint(v.Field(i).Uint(), 10)

		case reflect.Int:
			val = strconv.FormatInt(v.Field(i).Int(), 10)
		case reflect.Float64:
			val = strconv.FormatFloat(v.Field(i).Float(), 'f', 6, 64)
		}
		json = GetName(v.Field(i))
		name := v.Field(i).Type().Name()
		result[i] = name + "/" + json + "/" + val
	}

	return result
}

func GetName(field reflect.Value) string {
	rt := reflect.TypeOf(field)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get("Alloc"), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v != "" {
			return f.Name
		}
	}
	return ""
}

func main() {
	messages := make(chan []string)

	ctx, _ := context.WithCancel(context.Background())

	go RunGetMetrics(ctx, 2, messages)

	go RunSender(ctx, 10, messages)

	<-ctx.Done()
}

func RunSender(ctx context.Context, duration int, messages chan []string) error {
	ticker := time.NewTicker(time.Duration(duration) * time.Second) // создаём таймер
	for {
		select {
		case <-ticker.C:
			messagesFromChan := <-messages
			for _, mes := range messagesFromChan {
				Send(mes)
			}
			return nil
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
