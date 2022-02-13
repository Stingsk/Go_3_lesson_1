package httputil

import (
	"context"
	"errors"
	"fmt"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

func RunSender(ctx context.Context, duration int, messages *metrics.SensorData, wg *sync.WaitGroup) error {
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
			return errors.New("аварийное завершение")
		}
	}
}

func send(send string) {
	endpoint := "http://localhost:8080/update/" + send
	// конструируем HTTP-клиент
	client := resty.New()

	response, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(endpoint)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// печатаем код ответа
	logrus.Info("Запрос: ", send)
	logrus.Info("Статус-код ", response.StatusCode())
	// и печатаем его
	logrus.Info(string(response.Body()))
}
