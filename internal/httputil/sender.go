package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/storage"

	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const protocol string = "http://"

func RunSender(ctx context.Context, duration int, messages *metrics.SensorData, wg *sync.WaitGroup, host string) error {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(duration) * time.Second)
	for {
		select {
		case <-ticker.C:
			messagesFromChan := messages.Get()
			for _, mes := range messagesFromChan {
				sendPost(mes, host)
			}
		case <-ctx.Done():
			return errors.New("crash agent")
		}
	}
}

func send(send string, host string) {
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

func sendPost(metric storage.Metric, host string) {
	endpoint := protocol + host + "/update/"
	client := resty.New()

	m, err := json.Marshal(metric)
	if err != nil {
		logrus.Error(err)
	}
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(m).
		Post(endpoint)

	if err != nil {
		logrus.Error(err)
	}

	// печатаем код ответа
	logrus.Info("Send: ", metric)
	logrus.Info("Status code ", response.StatusCode())
	// и печатаем его
	logrus.Info(string(response.Body()))
}
