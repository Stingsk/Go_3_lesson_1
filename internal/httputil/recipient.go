package httputil

import (
	"context"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var metricData = make(map[storage.MetricName]storage.Metric)

func RunRecipient(wg *sync.WaitGroup, sigChan chan os.Signal) {
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
	apiRouter.Post("/update/*", recipientPost)
	apiRouter.Get("/value/*", recipientGet)
	apiRouter.Get("/", recipientGetAllMetrics)

	logrus.Info("Starting HTTP server")
	return apiRouter
}

func recipientPost(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	s := strings.Split(r.URL.Path, "/")
	if len(s) != 5 {
		http.Error(w, "Only 3 params in request are allowed!", http.StatusNotFound)
		return
	}

	var metric storage.Metric
	metric.NewMetricString(strings.ToLower(s[3]), strings.ToLower(s[2]), strings.ToLower(s[4]))
	if _, err := strconv.ParseFloat(s[4], 64); err != nil {
		http.Error(w, "Only Numbers  params in request are allowed!", http.StatusBadRequest)
		return
	}

	if metric.GetMetricType().IsZero() {
		http.Error(w, "MetricType NotImplemented!", http.StatusNotImplemented)
		return
	}
	if metric.GetMetricName().IsZero() {
		http.Error(w, "MetricName NotImplemented!", http.StatusNotImplemented)
		return
	}
	var valueMetric, found = metricData[metric.GetMetricName()]
	if found {
		logrus.Info("Данные обновлены")
		metricData[metric.GetMetricName()] = valueMetric.UpdateMetric(strings.ToLower(s[4]))
	} else {
		logrus.Info("Данныу добавлены")
		metricData[metric.GetMetricName()] = metric
	}

	logrus.Info(r.RequestURI)
}

func recipientGet(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	s := strings.Split(r.URL.Path, "/")
	if len(s) != 4 {
		http.Error(w, "Only 2 params in request are allowed!", http.StatusNotFound)
		return
	}

	var metric storage.Metric

	metric.NewMetricString(strings.ToLower(s[3]), strings.ToLower(s[2]), "")

	if metric.GetMetricType().IsZero() {
		http.Error(w, "MetricType NotImplemented!", http.StatusNotFound)
		return
	}
	if metric.GetMetricName().IsZero() {
		http.Error(w, "MetricName NotImplemented!", http.StatusNotFound)
		return
	}
	var valueMetric, found = metricData[metric.GetMetricName()]
	if found && valueMetric.GetMetricType().String() == strings.ToLower(s[2]) {
		logrus.Info("Данные обновлены")
		w.Write([]byte(valueMetric.GetValue()))
	} else {
		http.Error(w, "Value NotFound!", http.StatusNotFound)
		return
	}

	logrus.Info(r.RequestURI)
}

func recipientGetAllMetrics(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte(getAllMetrics()))

	logrus.Info(r.RequestURI)
}

func getAllMetrics() string {
	s := ""
	for _, element := range metricData {
		s += element.GetMetricName().String() + ": " + element.GetValue() + "\r"
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
