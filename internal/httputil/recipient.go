package httputil

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/file"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

var syncWrite = true

const gauge string = "gauge"
const counter string = "counter"

type MyMetric struct {
	Inner    *storage.MetricResourceMap
	FilePath string
}

func NewMyMetric(metric *storage.MetricResourceMap, filePath string) MyMetric {
	return MyMetric{
		Inner:    metric,
		FilePath: filePath,
	}
}

func RunServer(wg *sync.WaitGroup,
	sigChan chan os.Signal,
	host string,
	metrics *storage.MetricResourceMap,
	storeFile string,
	storeInterval time.Duration) {
	defer wg.Done()
	r := NewMyMetric(metrics, storeFile)
	server := &http.Server{Addr: getHost(host), Handler: service(&r)}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigChan
		file.WriteMetrics(storeFile, r.Inner.Metric)
		logrus.Info("Save data before Shutdown to " + storeFile)
		err := server.Shutdown(ctx)
		if err != nil {
			logrus.Fatal(err)
		}
		cancel()
	}()

	if storeInterval > 0 {
		go func() {
			ticker := time.NewTicker(storeInterval)
			for {
				<-ticker.C
				file.WriteMetrics(storeFile, r.Inner.Metric)
			}
		}()
		syncWrite = false
	}

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logrus.Fatal(err)
	}

	// Wait for server context to be stopped
	<-ctx.Done()
}

func service(metrics *MyMetric) http.Handler {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)
	apiRouter.Post("/update/", metrics.savePostMetric)
	apiRouter.Post("/value/", metrics.getValueMetric)
	apiRouter.Post("/update/{type}/{name}/{value}", metrics.saveMetric)
	apiRouter.Get("/value/{type}/{name}", metrics.getMetric)
	apiRouter.Get("/", metrics.getAllMetrics)
	apiRouter.Post("/update/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, getJSONError("Method NotFound!"), http.StatusNotFound)
	})

	logrus.Info("Starting HTTP server")

	return apiRouter
}

func (metrics *MyMetric) savePostMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logrus.Info("Url request: " + r.RequestURI)

	if r.Header.Get("Content-Type") != "application/json" {

		http.Error(w, getJSONError("Only application/json  can be Content-Type"), http.StatusUnsupportedMediaType)
	}
	defer r.Body.Close()

	var m storage.Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, getJSONError("Fail on parse request"), http.StatusBadRequest)
		return
	}

	if m.ID == "" || m.MType == "" {
		http.Error(w, getJSONError("ID or MType is empty"), http.StatusBadRequest)
	}

	var valueMetric, found = metrics.Inner.Metric[strings.ToLower(m.ID)]
	if found {
		if m.Delta != nil || m.Value != nil {
			updatedValueMetric := storage.NewMetricResource(storage.Update(m, valueMetric.Metric))
			metrics.Inner.Metric[strings.ToLower(m.ID)] = updatedValueMetric
			render.JSON(w, r, &updatedValueMetric.Metric)
			logrus.Info("Update data")
		} else {
			http.Error(w, getJSONError("Data is empty"), http.StatusBadRequest)
		}
	} else {
		if m.Delta != nil || m.Value != nil {
			metrics.Inner.Metric[strings.ToLower(m.ID)] = storage.NewMetricResource(m)
			render.JSON(w, r, &m)
			logrus.Info("Add data")
		} else {
			http.Error(w, getJSONError("Data is empty"), http.StatusBadRequest)
		}
	}

	if syncWrite {
		file.WriteMetrics(metrics.FilePath, metrics.Inner.Metric)
	}
	logrus.Info(r.RequestURI)
}

func (metrics *MyMetric) saveMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)

	metricType := strings.ToLower(chi.URLParam(r, "type"))
	if metricType != gauge && metricType != counter {
		http.Error(w, "Only metricType  gauge and counter in request are allowed!", http.StatusNotImplemented)
	}
	metricName := strings.ToLower(chi.URLParam(r, "name"))
	metricValue := strings.ToLower(chi.URLParam(r, "value"))

	_, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "Only Numbers  params in request are allowed!", http.StatusBadRequest)
		return
	}

	var valueMetric, found = metrics.Inner.Metric[metricName]
	if found {
		updatedValueMetric, err := storage.UpdateMetric(metricValue, valueMetric.Metric)
		if err != nil {
			http.Error(w, "Fail on update metric", http.StatusBadRequest)
			return
		}
		metrics.Inner.Metric[metricName] = storage.NewMetricResource(updatedValueMetric)
		logrus.Info("Updated data")
	} else {
		metric, _ := storage.NewMetric(metricValue, metricType, metricName)
		metrics.Inner.Metric[metricName] = storage.NewMetricResource(metric)
		logrus.Info("Added data")
	}

	if syncWrite {
		file.WriteMetrics(metrics.FilePath, metrics.Inner.Metric)
	}
	logrus.Info(r.RequestURI)
}

func (metrics *MyMetric) getValueMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logrus.Info("Url request: " + r.RequestURI)

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, getJSONError("Only application/json  can be Content-Type"), http.StatusUnsupportedMediaType)
	}
	defer r.Body.Close()

	var m storage.Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, getJSONError("Fail on parse request"), http.StatusBadRequest)
		return
	}

	if m.ID == "" || m.MType == "" {
		http.Error(w, getJSONError("ID or MType is empty"), http.StatusBadRequest)
	}

	var valueMetric, found = metrics.Inner.Metric[strings.ToLower(m.ID)]
	if found && m.Delta == nil && m.Value == nil {
		render.JSON(w, r, &valueMetric.Metric)
		logrus.Info("Send data")
	} else {
		http.Error(w, getJSONError("Data Not Found"), http.StatusNotFound)
	}

	logrus.Info(w.Header().Get("Content-Type"))
	logrus.Info(r.Body)
	logrus.Info(r.Header)
}

func (metrics *MyMetric) getMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)
	metricType := strings.ToLower(chi.URLParam(r, "type"))
	metricName := strings.ToLower(chi.URLParam(r, "name"))

	if metricType == "" {
		http.Error(w, "MetricType NotImplemented!", http.StatusNotFound)
		return
	}
	if metricName == "" {
		http.Error(w, "MetricName NotImplemented!", http.StatusNotFound)
		return
	}
	var valueMetric, found = metrics.Inner.Metric[metricName]
	if found && valueMetric.Metric.GetMetricType() == metricType {
		logrus.Info("Data received: " + valueMetric.Metric.GetValue())
		w.Write([]byte(valueMetric.Metric.GetValue()))
	} else {
		http.Error(w, "Value NotFound!", http.StatusNotFound)
		return
	}
}

func (metrics *MyMetric) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)
	metricsString := concatenationMetrics(metrics.Inner.Metric)
	logrus.Info("Data received: " + metricsString)
	w.Write([]byte(metricsString))

	logrus.Info(r.RequestURI)
}

func concatenationMetrics(metrics map[string]storage.MetricResource) string {
	s := ""
	for name, element := range metrics {
		s += name + ": " + element.Metric.GetValue() + "\r"
	}

	return s
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)

	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(middleware.NoCache)
	//router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.Timeout(60 * time.Second))
}

func getJSONError(errorText string) string {
	return "{ \"error\" : \"" + errorText + "\"}"
}

func getHost(host string) string {
	re := regexp.MustCompile(":[0-9]+")
	port := re.FindAllString(host, 1)
	return port[0]
}
