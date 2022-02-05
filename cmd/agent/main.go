package main

import (
	"bytes"
	"fmt"
	metrics "github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

func main() {
	monitor, err := metrics.NewMonitor()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	v := reflect.ValueOf(monitor)

	values := make(map[string]string)

	for i := 0; i < v.NumField(); i++ {
		fiald := ""
		switch v.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fiald = strconv.FormatInt(v.Field(i).Int(), 10)
		case reflect.String:
			fiald = v.Field(i).String()
		case reflect.Uint64, reflect.Uint32:
			fiald = strconv.FormatUint(v.Field(i).Uint(), 10)
		case reflect.Float64:
			fiald = strconv.FormatFloat(v.Field(i).Float(), 'f', 6, 64)
		}
		values[fiald] = v.Field(i).Type().String()
	}

	for val, t := range values {
		Send(val, t)
	}
}

func Send(v string, t string) {
	endpoint := "http://localhost:8080/update/" + v + "/" + t + "/"
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
