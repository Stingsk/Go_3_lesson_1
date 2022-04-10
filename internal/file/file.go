package file

import (
	"encoding/json"
	"errors"
	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
	"io/ioutil"
	"os"
)

func WriteMetrics(fileName string, metricResourceMap *storage.MetricResourceMap) error {
	if fileName == "" {
		return errors.New("empty filename")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	if metricResourceMap == nil {
		return nil
	}
	metricResourceMap.Mutex.Lock()
	defer metricResourceMap.Mutex.Unlock()
	marshalMetric, err := json.Marshal(metricResourceMap.Metric)
	if err != nil {
		return err
	}
	file.Write(marshalMetric)
	file.Close()

	return nil
}

func ReadMetrics(fileName string) (map[string]storage.Metric, error) {
	if fileName == "" {
		return nil, nil
	}
	metricData := make(map[string]storage.Metric)
	fileRead, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return metricData, err
	}

	defer fileRead.Close()

	data, err := ioutil.ReadAll(fileRead)
	if err != nil {
		return metricData, err
	}
	err = json.Unmarshal(data, &metricData)

	if err != nil {
		return metricData, err
	}

	return metricData, nil
}
