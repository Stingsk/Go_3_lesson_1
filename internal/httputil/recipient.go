package httputil

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

var metricData = make(map[string]storage.Metric)

const gauge string = "gauge"
const counter string = "counter"
const host string = "localhost:8080"

func RunServer(wg *sync.WaitGroup, sigChan chan os.Signal) {
	defer wg.Done()
	server := &http.Server{Addr: host, Handler: service()}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigChan
		err := server.Shutdown(ctx)
		if err != nil {
			logrus.Fatal(err)
		}

		cancel()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logrus.Fatal(err)
	}

	// Wait for server context to be stopped
	<-ctx.Done()
}

func service() http.Handler {

	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)
	apiRouter.Post("/update/"+gauge+"/{gauge}/{value}", postGaugeMetric)
	apiRouter.Post("/update/"+counter+"/{counter}/{value}", postCounterMetric)
	apiRouter.Get("/value/{type}/{name}", getMetric)
	apiRouter.Get("/", getAllMetrics)
	apiRouter.Post("/update/"+gauge+"*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method NotImplemented!", http.StatusNotFound)
		return
	})
	apiRouter.Post("/update/"+counter+"*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method NotImplemented!", http.StatusNotFound)
		return
	})
	apiRouter.Post("/update/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method NotImplemented!", http.StatusNotImplemented)
		return
	})

	logrus.Info("Starting HTTP server")
	return apiRouter
}

func postGaugeMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)

	metricName := chi.URLParam(r, gauge)
	metricValue := chi.URLParam(r, "value")

	var metricNameType storage.MetricName
	metricNameType.NewMetricName(metricName)

	if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
		http.Error(w, "Only Numbers  params in request are allowed!", http.StatusBadRequest)
		return
	}

	if metricNameType.IsZero() {
		http.Error(w, "MetricName NotImplemented!", http.StatusNotImplemented)
		return
	}

	var valueMetric, found = metricData[metricName]
	if found {
		valueMetric.UpdateMetric(metricValue, gauge)
		metricData[metricName] = valueMetric
		logrus.Info("Updated data")
	} else {
		var metric storage.Metric
		metric.NewMetric(strings.ToLower(metricName), strings.ToLower(gauge), strings.ToLower(metricValue))
		metricData[metricName] = metric
		logrus.Info("Added data")
	}

	logrus.Info(r.RequestURI)
}

func postCounterMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)

	metricName := chi.URLParam(r, counter)
	metricValue := chi.URLParam(r, "value")

	var metricNameType storage.MetricName
	metricNameType.NewMetricName(metricName)

	if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
		http.Error(w, "Only Numbers  params in request are allowed!", http.StatusBadRequest)
		return
	}

	var valueMetric, found = metricData[metricName]
	if found {
		valueMetric.UpdateMetric(metricValue, counter)
		metricData[metricName] = valueMetric
		logrus.Info("Updated data")
	} else {
		var metric storage.Metric
		metric.NewMetric(strings.ToLower(metricName), strings.ToLower(counter), strings.ToLower(metricValue))
		metricData[metricName] = metric
		logrus.Info("Added data")
	}

	logrus.Info(r.RequestURI)
}

func getMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if metricType == "" {
		http.Error(w, "MetricType NotImplemented!", http.StatusNotFound)
		return
	}
	if metricName == "" {
		http.Error(w, "MetricName NotImplemented!", http.StatusNotFound)
		return
	}
	var valueMetric, found = metricData[metricName]
	if found && valueMetric.GetMetricType() == metricType {
		logrus.Info("Data received: " + valueMetric.GetValue())
		w.Write([]byte(valueMetric.GetValue()))
	} else {
		http.Error(w, "Value NotFound!", http.StatusNotFound)
		return
	}
}

func getAllMetrics(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)
	metrics := concatenationMetrics()
	logrus.Info("Data received: " + metrics)
	w.Write([]byte(metrics))

	logrus.Info(r.RequestURI)
}

func concatenationMetrics() string {
	s := ""
	for _, element := range metricData {
		s += element.GetMetricName() + ": " + element.GetValue() + "\r"
	}

	return s
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)

	router.Use(
		middleware.SetHeader("Content-Type", "text/plain"),
	)
	router.Use(middleware.NoCache)
	router.Use(middleware.AllowContentType("text/plain"))
	router.Use(middleware.Timeout(60 * time.Second))
}
