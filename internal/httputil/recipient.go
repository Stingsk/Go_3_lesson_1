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

func RunServer(wg *sync.WaitGroup, sigChan chan os.Signal, host string) {
	defer wg.Done()
	server := &http.Server{Addr: getHost(host), Handler: service()}
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
	apiRouter.Post("/update/", postJSONMetric)
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
func postJSONMetric(w http.ResponseWriter, r *http.Request) {
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

	var valueMetric, found = metricData[strings.ToLower(m.ID)]
	if found {
		if m.Delta != nil || m.Value != nil {
			updatedValueMetric := storage.Update(m, valueMetric)
			metricData[strings.ToLower(m.ID)] = updatedValueMetric
			render.JSON(w, r, &updatedValueMetric)
			logrus.Info("Update data")
		} else {
			http.Error(w, getJSONError("Data is empty"), http.StatusBadRequest)
		}
	} else {
		if m.Delta != nil || m.Value != nil {
			metricData[strings.ToLower(m.ID)] = m
			render.JSON(w, r, &m)
			logrus.Info("Add data")
		} else {
			http.Error(w, getJSONError("Data is empty"), http.StatusBadRequest)
		}
	}

	logrus.Info(w.Header().Get("Content-Type"))
	logrus.Info(r.Body)
	logrus.Info(r.Header)
}

func postValueMetric(w http.ResponseWriter, r *http.Request) {
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

	var valueMetric, found = metricData[strings.ToLower(m.ID)]
	if found && m.Delta == nil && m.Value == nil {
		render.JSON(w, r, &valueMetric)
		logrus.Info("Send data")
	} else {
		http.Error(w, getJSONError("Data Not Found"), http.StatusNotFound)
	}

	logrus.Info(w.Header().Get("Content-Type"))
	logrus.Info(r.Body)
	logrus.Info(r.Header)
}

func postGaugeMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)

	metricName := strings.ToLower(chi.URLParam(r, gauge))
	metricValue := strings.ToLower(chi.URLParam(r, "value"))

	_, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "Only Numbers  params in request are allowed!", http.StatusBadRequest)
		return
	}

	var valueMetric, found = metricData[metricName]
	if found {
		updatedValueMetric, err := storage.UpdateMetric(metricValue, valueMetric)
		if err != nil {
			http.Error(w, "Fail on update metric", http.StatusBadRequest)
			return
		}
		metricData[metricName] = updatedValueMetric
		logrus.Info("Updated data")
	} else {
		metric, _ := storage.NewMetric(metricValue, gauge, metricName)
		metricData[metricName] = metric
		logrus.Info("Added data")
	}

	logrus.Info(r.RequestURI)
}

func postCounterMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)

	metricName := strings.ToLower(chi.URLParam(r, counter))
	metricValue := strings.ToLower(chi.URLParam(r, "value"))

	if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
		http.Error(w, "Only Numbers  params in request are allowed!", http.StatusBadRequest)
		return
	}

	var valueMetric, found = metricData[metricName]
	if found {
		updatedValueMetric, err := storage.UpdateMetric(metricValue, valueMetric)
		if err != nil {
			http.Error(w, "Fail on update metric", http.StatusBadRequest)
			return
		}
		metricData[metricName] = updatedValueMetric
		logrus.Info("Updated data")
	} else {
		metric, err := storage.NewMetric(metricValue, counter, metricName)
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

	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(middleware.NoCache)
	//router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.Timeout(60 * time.Second))
}

func getJSONError(errorText string) string {
	return "{ \"error\" : \"" + errorText + "\"}"
}

func getHost(host string) string {
	re := regexp.MustCompile("[0-9]+")
	port := re.FindAllString(host, 1)
	return ":" + port[0]
}
