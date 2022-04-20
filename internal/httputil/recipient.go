package httputil

import (
	"compress/gzip"
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

type gzipResponseWriter struct {
	Writer *gzip.Writer
	http.ResponseWriter
}
type ServerConfig struct {
	WaitGroup     *sync.WaitGroup
	SigChan       chan os.Signal
	Host          string
	Metrics       map[string]storage.Metric
	StoreFile     string
	StoreInterval time.Duration
	SignKey       string
}

var MetricLocal *storage.MetricResourceMap
var StoreFile string
var SignKey string

func RunServer(serverConfig ServerConfig) {
	StoreFile = serverConfig.StoreFile
	SignKey = serverConfig.SignKey
	MetricLocal = &storage.MetricResourceMap{
		Metric:     nil,
		Mutex:      sync.Mutex{},
		Repository: &storage.MemoryStorage{},
	}
	MetricLocal.Metric = serverConfig.Metrics
	defer serverConfig.WaitGroup.Done()
	server := &http.Server{Addr: getHost(serverConfig.Host), Handler: service()}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-serverConfig.SigChan
		file.WriteMetrics(serverConfig.StoreFile, MetricLocal)
		logrus.Info("Save data before Shutdown to " + serverConfig.StoreFile)
		err := server.Shutdown(ctx)
		if err != nil {
			logrus.Fatal(err)
		}
		cancel()
	}()

	if serverConfig.StoreInterval > 0 {
		go func() {
			ticker := time.NewTicker(serverConfig.StoreInterval)
			for {
				<-ticker.C
				file.WriteMetrics(serverConfig.StoreFile, MetricLocal)
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

	if !m.IsHashValid(SignKey) {
		http.Error(w, getJSONError("Hash is invalid"), http.StatusBadRequest)
	}

	var valueMetric, err = MetricLocal.Repository.UpdateMetric(MetricLocal, m)
	if err != nil {
		http.Error(w, getJSONError("Data is empty"), http.StatusBadRequest)
	}

	render.JSON(w, r, valueMetric)

	if syncWrite {
		file.WriteMetrics(StoreFile, MetricLocal)
	}
	logrus.Info(r.RequestURI)
}

func saveMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	logrus.Info("Url request: " + r.RequestURI)

	metricType := strings.ToLower(chi.URLParam(r, "type"))
	if metricType != gauge && metricType != counter {
		http.Error(w, "Only metricType  gauge and counter in request are allowed!", http.StatusNotImplemented)
	}
	metricName := strings.ToLower(chi.URLParam(r, "name"))
	metricValue := strings.ToLower(chi.URLParam(r, "value"))

	var _, err = MetricLocal.Repository.UpdateMetricByParameters(MetricLocal, metricName, metricType, metricValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if syncWrite {
		file.WriteMetrics(StoreFile, MetricLocal)
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

	var valueMetric, found = MetricLocal.Metric[strings.ToLower(m.ID)]
	m.SetHash(SignKey)
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
	w.Header().Set("Content-Type", "text/plain")
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
	var valueMetric, found = MetricLocal.Metric[metricName]
	if found && valueMetric.GetMetricType() == metricType {
		logrus.Info("Data received: " + valueMetric.GetValue())
		w.Write([]byte(valueMetric.GetValue()))
	} else {
		http.Error(w, "Value NotFound!", http.StatusNotFound)
		return
	}
}

func getAllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	logrus.Info("Url request: " + r.RequestURI)
	metricsString := concatenationMetrics(MetricLocal.Metric)
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
	router.Use(middleware.NoCache)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(gzipResponse)
}

func getJSONError(errorText string) string {
	return "{ \"error\" : \"" + errorText + "\"}"
}

func getHost(host string) string {
	re := regexp.MustCompile(":[0-9]+")
	port := re.FindAllString(host, 1)
	return port[0]
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func gzipResponse(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		next.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	}

	return http.HandlerFunc(fn)
}
