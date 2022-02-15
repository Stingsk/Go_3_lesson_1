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

func RunServer(wg *sync.WaitGroup, sigChan chan os.Signal) {
	defer wg.Done()
	server := &http.Server{Addr: "localhost:8080", Handler: service()}
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
	apiRouter.Post("/update/*", postMetric)
	apiRouter.Get("/value/*", getMetric)
	apiRouter.Get("/", getAllMetrics)

	logrus.Info("Starting HTTP server")
	return apiRouter
}

func postMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)
	s := strings.Split(r.URL.Path, "/")
	if len(s) != 5 {
		http.Error(w, "Only 3 params in request are allowed!", http.StatusNotFound)
		return
	}

	metricType := strings.ToLower(s[2])
	metricName := strings.ToLower(s[3])
	metricValue := strings.ToLower(s[4])

	var metricNameType storage.MetricName
	metricNameType.NewMetricName(metricName)

	if _, err := strconv.ParseFloat(s[4], 64); err != nil {
		http.Error(w, "Only Numbers  params in request are allowed!", http.StatusBadRequest)
		return
	}

	if metricType == "gauge" && metricNameType.IsZero() {
		http.Error(w, "MetricName NotImplemented!", http.StatusNotImplemented)
		return
	}
	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "MetricName NotImplemented!", http.StatusNotImplemented)
		return
	}

	var valueMetric, found = metricData[metricName]
	if found {
		valueMetric.UpdateMetric(metricValue, metricType)
		metricData[metricName] = valueMetric
		logrus.Info("Updated data")
	} else {
		var metric storage.Metric
		metric.NewMetric(strings.ToLower(s[3]), strings.ToLower(s[2]), strings.ToLower(s[4]))
		metricData[metricName] = metric
		logrus.Info("Added data")
	}

	logrus.Info(r.RequestURI)
}

func getMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)
	s := strings.Split(r.URL.Path, "/")
	if len(s) != 4 {
		http.Error(w, "Only 2 params in request are allowed!", http.StatusNotFound)
		return
	}

	metricType := strings.ToLower(s[2])
	metricName := strings.ToLower(s[3])

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
