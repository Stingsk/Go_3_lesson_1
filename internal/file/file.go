package file

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
	"github.com/sirupsen/logrus"
)

type Event struct {
	ID     string         `json:"Id"`
	Metric storage.Metric `json:"Metric"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}
func (p *producer) WriteEvent(event *Event) error {
	return p.encoder.Encode(&event)
}
func (p *producer) Close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}
func (c *consumer) ReadEvent() (*Event, error) {
	event := &Event{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}
func (c *consumer) Close() error {
	return c.file.Close()
}

func WriteMetrics(fileName string, events map[string]storage.Metric) {
	if fileName == "" {
		return
	}
	producer, err := NewProducer(fileName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer producer.Close()
	for _, event := range events {
		eventToWrite := Event{
			ID:     event.ID,
			Metric: event,
		}
		if err := producer.WriteEvent(&eventToWrite); err != nil {
			logrus.Fatal(err)
		}
	}
}

func ReadMetrics(fileName string) (map[string]storage.Metric, error) {
	if fileName == "" {
		return nil, nil
	}
	var metricData = make(map[string]storage.Metric)
	fileRead, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}

	reader := bufio.NewReader(fileRead)
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 100*1024*1024), 100*1024*1024) // читаем большие файлы
	defer fileRead.Close()

	for scanner.Scan() {
		event := &Event{}
		json.Unmarshal(scanner.Bytes(), &event)
		metricData[event.ID] = event.Metric
	}

	return metricData, nil
}
