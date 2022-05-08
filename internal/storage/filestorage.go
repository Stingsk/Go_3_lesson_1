package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

type FileStorage struct {
	filePath    string
	syncChannel chan struct{}
	metrics     map[string]*Metric
	Mutex       sync.Mutex
}

func NewFileStorage(filePath string, syncChannel chan struct{}) (*FileStorage, error) {
	var fs FileStorage

	metricsCache := make(map[string]*Metric)
	fs = FileStorage{
		filePath:    filePath,
		syncChannel: syncChannel,
		metrics:     metricsCache,
	}

	return &fs, nil
}

func (fs *FileStorage) UpdateMetric(_ context.Context, metric Metric) error {
	fs.Mutex.Lock()
	defer fs.sync()
	defer fs.Mutex.Unlock()
	var valueMetric = fs.metrics[metric.ID]
	if valueMetric.GetValue() != "" {
		if metric.Delta != nil || metric.Value != nil {
			valueMetric.Update(metric)
			fs.metrics[metric.ID] = valueMetric
			return nil
		} else {
			return errors.New("data is empty")
		}
	} else {
		if metric.Delta != nil || metric.Value != nil {
			fs.metrics[metric.ID] = &metric
			return nil
		} else {
			return errors.New("data is empty")
		}
	}
}

func (fs *FileStorage) UpdateMetricByParameters(_ context.Context, metricName string, metricType string, value string) error {
	fs.Mutex.Lock()
	defer fs.sync()
	defer fs.Mutex.Unlock()
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return errors.New("only numbers  params in request are allowed")
	}

	var valueMetric, found = fs.metrics[metricName]
	if found {
		err := valueMetric.UpdateMetricResource(value)
		fs.metrics[metricName] = valueMetric
		if err != nil {
			return err
		}
		return nil
	} else {
		err := fs.newMetric(value, metricType, metricName)
		if err != nil {
			return err
		}

		return nil
	}
}

func (fs *FileStorage) GetMetric(_ context.Context, name string, _ string) (*Metric, error) {
	if name == "" {
		return &Metric{}, errors.New("empty name")
	}
	value, ok := fs.metrics[name]

	if !ok {
		return &Metric{}, errors.New("name not found")
	}

	return value, nil
}
func (fs *FileStorage) GetMetrics(_ context.Context) (map[string]*Metric, error) {
	return fs.metrics, nil
}

func (fs *FileStorage) UpdateMetrics(_ context.Context, metricsBatch []*Metric) error {
	fs.Mutex.Lock()
	defer fs.sync()
	defer fs.Mutex.Unlock()
	fs.metrics = make(map[string]*Metric)
	for _, mb := range metricsBatch {
		fs.metrics[mb.ID] = mb
	}

	return nil
}

func (fs *FileStorage) WriteMetrics() error {
	if fs.filePath == "" {
		return errors.New("empty filename")
	}

	file, err := os.Create(fs.filePath)
	if err != nil {
		return err
	}
	if fs.metrics == nil {
		return nil
	}
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()
	marshalMetric, err := json.Marshal(fs.metrics)
	if err != nil {
		return err
	}
	file.Write(marshalMetric)
	file.Close()

	return nil
}

func (fs *FileStorage) ReadMetrics() error {
	if fs.filePath == "" {
		return errors.New("empty filename")
	}
	metricData := make(map[string]*Metric)
	fileRead, err := os.OpenFile(fs.filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer fileRead.Close()
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	data, err := ioutil.ReadAll(fileRead)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &metricData)

	if err != nil {
		return err
	}

	fs.metrics = metricData

	return nil
}

func (fs *FileStorage) Ping(_ context.Context) error {
	if fs.filePath == "" {
		return errors.New("empty filename")
	}
	fileRead, err := os.OpenFile(fs.filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer fileRead.Close()
	_, err = fileRead.Stat()

	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) sync() {
	//fs.syncChannel <- struct{}{}
}

func (fs *FileStorage) newMetric(value string, metricType string, name string) error {
	metric := Metric{
		ID:    name,
		MType: strings.ToLower(metricType),
		Delta: nil,
		Value: nil,
	}
	if strings.ToLower(metricType) == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		metric.Value = &v
	} else if strings.ToLower(metric.MType) == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		delta := int64(0)
		if metric.Delta != nil {
			delta = int64(*metric.Delta)
		}
		metric.Delta = sumInt(delta, newValue)
	}
	fs.metrics[name] = &metric

	return nil
}
