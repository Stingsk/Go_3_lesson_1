package httputil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func RunSender(ctx context.Context, duration int, messages *metrics.SensorData, wg *sync.WaitGroup, sigChan chan os.Signal) error {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(duration) * time.Second) // создаём таймер
	for {
		select {
		case <-ticker.C:
			messagesFromChan := messages.Get()
			for _, mes := range messagesFromChan {
				send(mes)
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-sigChan:
			return errors.New("аварийное завершение")
		}
	}
}

func send(send string) {
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
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len("monitor")))
	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// печатаем код ответа
	fmt.Println("Статус-код ", response.Status)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// и печатаем его
	fmt.Println(string(body))
}
