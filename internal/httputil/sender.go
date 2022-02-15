package httputil

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const protocol string = "http://"

func RunSender(ctx context.Context, duration int, messages *metrics.SensorData, wg *sync.WaitGroup) error {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(duration) * time.Second)
	for {
		select {
		case <-ticker.C:
			messagesFromChan := messages.Get()
			for _, mes := range messagesFromChan {
				send(mes)
			}
		case <-ctx.Done():
			return errors.New("crash agent")
		}
	}
}

func send(send string) {
	endpoint := protocol + host + "/update/" + send
	client := resty.New()

	response, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(endpoint)

	if err != nil {
		fmt.Println(err)
	}

	// печатаем код ответа
	logrus.Info("Send: ", send)
	logrus.Info("Status code ", response.StatusCode())
	// и печатаем его
	logrus.Info(string(response.Body()))
}
