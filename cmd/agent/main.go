package main

import (
	"bytes"
	"fmt"
	metrics "github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {

	endpoint := "http://localhost:8080/setMetrics/"

	monitor, err := metrics.NewMonitor()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// конструируем HTTP-клиент
	client := &http.Client{}
	// конструируем запрос
	// запрос методом POST должен, кроме заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	// в большинстве случаев отлично подходит bytes.Buffer
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(monitor))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Content-Length", strconv.Itoa(len(monitor)))
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
