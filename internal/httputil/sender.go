package httputil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/storage"

	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const protocol string = "http://"

func RunSender(ctx context.Context, duration time.Duration, messages *metrics.SensorData, wg *sync.WaitGroup, host string) {
	defer wg.Done()
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			messagesFromChan := messages.Get()
			for _, mes := range messagesFromChan {
				sendPost(*mes.Metric, host)
			}
		case <-ctx.Done():
			logrus.Error("crash agent")
			return
		}
	}
}

func send(send string, host string) {
	endpoint := getHostSend(host) + "/update/" + send
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
	endpoint := getHostSend(host) + "/update/"
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

func getHostSend(host string) string {
	if strings.Contains(host, protocol) {
		return host
	}

	return protocol + host
}
