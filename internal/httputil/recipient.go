package httputil

import (
	"context"
	"encoding/json"
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
	"github.com/go-chi/render"
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
	apiRouter.Post("/update/", postJsonMetric)
	apiRouter.Post("/value/", postValueMetric)
	apiRouter.Post("/update/"+gauge+"/{gauge}/{value}", postGaugeMetric)
	apiRouter.Post("/update/"+counter+"/{counter}/{value}", postCounterMetric)
	apiRouter.Get("/value/{type}/{name}", getMetric)
	apiRouter.Get("/", getAllMetrics)
	apiRouter.Post("/update/"+gauge+"*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method NotImplemented!", http.StatusNotFound)
	})
	apiRouter.Post("/update/"+counter+"*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method NotImplemented!", http.StatusNotFound)
	})
	apiRouter.Post("/update/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method NotImplemented!", http.StatusNotImplemented)
	})

	logrus.Info("Starting HTTP server")

	return apiRouter
}
func postJsonMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logrus.Info("Url request: " + r.RequestURI)

	if r.Header.Get("Content-Type") != "application/json" {

		http.Error(w, "Only application/json  can be Content-Type", http.StatusUnsupportedMediaType)
	}
	defer r.Body.Close()

	var m storage.Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, "Fail on parse request", http.StatusBadRequest)
		return
	}

	var valueMetric, found = metricData[m.ID]
	if found {
		if m.Delta != nil || m.Value != nil {
			updatedValueMetric := storage.Update(m, valueMetric)
			metricData[m.ID] = updatedValueMetric
			logrus.Info("Update data")
		} else {
			http.Error(w, "{}", http.StatusBadRequest)
		}
	} else {
		metricData[m.ID] = m
		logrus.Info("Add data")
	}

	logrus.Info(w.Header().Get("Content-Type"))
	logrus.Info(r.Body)
	logrus.Info(r.Header)
	return
}

func postValueMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logrus.Info("Url request: " + r.RequestURI)

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Only application/json  can be Content-Type", http.StatusUnsupportedMediaType)
	}
	defer r.Body.Close()

	var m storage.Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	var valueMetric, found = metricData[m.ID]
	if found {
		if m.Delta == nil && m.Value == nil {
			render.JSON(w, r, &valueMetric)
			logrus.Info("Send data")
		} else {
			http.Error(w, "{}", http.StatusNotFound)
		}
	} else {
		http.Error(w, "{}", http.StatusNotFound)
	}

	logrus.Info(w.Header().Get("Content-Type"))
	logrus.Info(r.Body)
	logrus.Info(r.Header)
	return
}

func postGaugeMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)

	metricName := chi.URLParam(r, gauge)
	metricValue := chi.URLParam(r, "value")

	_, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "Only Numbers  params in request are allowed!", http.StatusBadRequest)
		return
	}

	var valueMetric, found = metricData[metricName]
	if found {
		updatedValueMetric, err := storage.UpdateMetric(strings.ToLower(metricValue), valueMetric)
		if err != nil {
			http.Error(w, "Fail on update metric", http.StatusBadRequest)
			return
		}
		metricData[metricName] = updatedValueMetric
		logrus.Info("Updated data")
	} else {
		metric, _ := storage.NewMetric(strings.ToLower(metricValue), strings.ToLower(gauge), metricName)
		metricData[metricName] = metric
		logrus.Info("Added data")
	}

	logrus.Info(r.RequestURI)
	return
}

func postCounterMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)

	metricName := chi.URLParam(r, counter)
	metricValue := chi.URLParam(r, "value")

	if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
		http.Error(w, "Only Numbers  params in request are allowed!", http.StatusBadRequest)
		return
	}

	var valueMetric, found = metricData[metricName]
	if found {
		updatedValueMetric, err := storage.UpdateMetric(strings.ToLower(metricValue), valueMetric)
		if err != nil {
			http.Error(w, "Fail on update metric", http.StatusBadRequest)
			return
		}
		metricData[metricName] = updatedValueMetric
		logrus.Info("Updated data")
	} else {
		metric, err := storage.NewMetric(strings.ToLower(metricValue), strings.ToLower(counter), metricName)
		if err != nil {
			http.Error(w, "Fail on add new metric", http.StatusBadRequest)
			return
		}
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
	for name, element := range metricData {
		s += name + ": " + element.GetValue() + "\r"
	}

	return s
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)

	router.Use(
		middleware.SetHeader("Content-Type", "application/json"),
	)
	router.Use(middleware.NoCache)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.Timeout(60 * time.Second))
}
