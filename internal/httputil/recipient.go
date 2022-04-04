package httputil

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
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

var metricLocal *storage.MetricResourceMap
var StoreFile string

func RunServer(wg *sync.WaitGroup,
	sigChan chan os.Signal,
	host string,
	metrics map[string]storage.Metric,
	storeFile string,
	storeInterval time.Duration) {
	StoreFile = storeFile
	metricLocal.Metric = metrics
	defer wg.Done()
	server := &http.Server{Addr: getHost(host), Handler: service()}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigChan
		file.WriteMetrics(storeFile, metricLocal)
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
				file.WriteMetrics(storeFile, metricLocal)
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

func service() http.Handler {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)
	apiRouter.Post("/update/", savePostMetric)
	apiRouter.Post("/value/", getValueMetric)
	apiRouter.Post("/update/{type}/{name}/{value}", saveMetric)
	apiRouter.Get("/value/{type}/{name}", getMetric)
	apiRouter.Get("/", getAllMetrics)
	apiRouter.Post("/update/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, getJSONError("Method NotFound!"), http.StatusNotFound)
	})

	logrus.Info("Starting HTTP server")

	return apiRouter
}

func savePostMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logrus.Info("Url request: " + r.RequestURI)

	if r.Header.Get("Content-Type") != "application/json" {

		http.Error(w, getJSONError("Only application/json  can be Content-Type"), http.StatusUnsupportedMediaType)
	}
	defer r.Body.Close()

	var m storage.Metric
	errDec := json.NewDecoder(r.Body).Decode(&m)
	if errDec != nil {
		http.Error(w, getJSONError("Fail on parse request"), http.StatusBadRequest)
		return
	}

	if m.ID == "" || m.MType == "" {
		http.Error(w, getJSONError("ID or MType is empty"), http.StatusBadRequest)
	}

	var valueMetric, err = storage.UpdateMetric(metricLocal, m)
	if err != nil {
		http.Error(w, getJSONError("Data is empty"), http.StatusBadRequest)
	}

	render.JSON(w, r, valueMetric)

	if syncWrite {
		file.WriteMetrics(StoreFile, metricLocal)
	}
	logrus.Info(r.RequestURI)
}

func saveMetric(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Url request: " + r.RequestURI)

	metricType := strings.ToLower(chi.URLParam(r, "type"))
	if metricType != gauge && metricType != counter {
		http.Error(w, "Only metricType  gauge and counter in request are allowed!", http.StatusNotImplemented)
	}
	metricName := strings.ToLower(chi.URLParam(r, "name"))
	metricValue := strings.ToLower(chi.URLParam(r, "value"))

	var _, err = storage.UpdateMetricByParameters(metricLocal, metricName, metricType, metricValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if syncWrite {
		file.WriteMetrics(StoreFile, metricLocal)
	}
	logrus.Info(r.RequestURI)
}

func getValueMetric(w http.ResponseWriter, r *http.Request) {
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

	var valueMetric, found = metricLocal.Metric[strings.ToLower(m.ID)]
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
	var valueMetric, found = metricLocal.Metric[metricName]
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
	metricsString := concatenationMetrics(metricLocal.Metric)
	logrus.Info("Data received: " + metricsString)
	w.Write([]byte(metricsString))

	logrus.Info(r.RequestURI)
}

func concatenationMetrics(metrics map[string]storage.Metric) string {
	s := ""
	for name, element := range metrics {
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
	re := regexp.MustCompile(":[0-9]+")
	port := re.FindAllString(host, 1)
	return port[0]
}
