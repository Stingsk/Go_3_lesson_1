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

	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Address            string        `env:"ADDRESS"`
	StoreInterval      time.Duration `env:"STORE_INTERVAL"`
	StoreFile          string        `env:"STORE_FILE"`
	Restore            bool          `env:"RESTORE"`
	SignKey            string        `env:"KEY"`
	DataBaseConnection string        `env:"DATABASE_DSN"`
	LogLevel           string        `env:"LogLevel"`
	WaitGroup          *sync.WaitGroup
	SigChan            chan os.Signal
}

const (
	requestTimeout = 1 * time.Second
)

const gauge string = "gauge"
const counter string = "counter"

type gzipResponseWriter struct {
	Writer *gzip.Writer
	http.ResponseWriter
}

var SignKey string
var Storage storage.Repository

func RunServer(serverConfig Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	SignKey = serverConfig.SignKey

	defer serverConfig.WaitGroup.Done()
	server := &http.Server{Addr: getHost(serverConfig.Address), Handler: service()}

	if serverConfig.DataBaseConnection != "" {
		logrus.Info("Start DataBase Store ")
		DBStore, err := storage.NewDBStore(serverConfig.DataBaseConnection)
		if err != nil {
			logrus.Info(err)
		}
		Storage = DBStore
	} else if serverConfig.StoreFile != "" {
		logrus.Info("Start File Store ")
		syncChannel := make(chan struct{}, 1)
		FileStorage, err := storage.NewFileStorage(serverConfig.StoreFile, syncChannel)
		if err != nil {
			logrus.Info(err)
		}

		if serverConfig.Restore {
			logrus.Info("Load data from  file")
			err := FileStorage.ReadMetrics()
			if err != nil {
				logrus.Info("Fail to restore data")
			}
		}

		Storage = FileStorage

		go func() {
			ticker := new(time.Ticker)
			if serverConfig.StoreInterval > 0 {
				ticker = time.NewTicker(serverConfig.StoreInterval)
			}
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					FileStorage.WriteMetrics()
				case <-syncChannel:
					if serverConfig.StoreInterval == 0 {
						FileStorage.WriteMetrics()
					}
				case <-serverConfig.SigChan:
					FileStorage.WriteMetrics()
					logrus.Info("Save data before Shutdown to " + serverConfig.StoreFile)
					err := server.Shutdown(ctx)
					if err != nil {
						logrus.Fatal(err)
					}
					cancel()
					return
				}
			}
		}()
	} else {
		logrus.Info("Start Memory Store ")
		MemoryStorage := storage.NewMemoryStorage()
		Storage = MemoryStorage
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
	apiRouter.Post("/updates/", savePostMetrics)
	apiRouter.Post("/value/", getValueMetric)
	apiRouter.Post("/update/{type}/{name}/{value}", saveMetric)
	apiRouter.Get("/value/{type}/{name}", getMetric)
	apiRouter.Get("/ping", pingDataBase)
	apiRouter.Get("/", getAllMetrics)
	apiRouter.Post("/update/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, getJSONError("Method NotFound!"), http.StatusNotFound)
	})

	logrus.Info("Starting HTTP server")

	return apiRouter
}

func savePostMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, getJSONError("Only application/json  can be Content-Type"), http.StatusUnsupportedMediaType)
		return
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
		return
	}

	if !m.IsHashValid(SignKey) {
		http.Error(w, getJSONError("Hash is invalid"), http.StatusBadRequest)
		return
	}

	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()
	err := Storage.UpdateMetric(requestContext, m)
	if err != nil {
		logrus.Error(err)
		http.Error(w, getJSONError("Data is empty"), http.StatusBadRequest)
		return
	}
	mar, _ := json.Marshal(m)
	logrus.Info("SavePostMetric Value: ", mar)
	render.JSON(w, r, m)

}

func savePostMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, getJSONError("Only application/json  can be Content-Type"), http.StatusUnsupportedMediaType)
		return
	}
	defer r.Body.Close()

	var m []*storage.Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		logrus.Error(err)
		http.Error(w, getJSONError("Fail on parse request"), http.StatusBadRequest)
		return
	}

	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()
	err = Storage.UpdateMetrics(requestContext, m)
	if err != nil {
		logrus.Error(err)
		http.Error(w, getJSONError("Data is empty"), http.StatusBadRequest)
		return
	}

	mar, _ := json.Marshal(m)
	logrus.Info("SavePostMetrics Value: ", mar)
	w.Write([]byte("{ \"success\" : \"success\"}"))
}

func saveMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	metricType := strings.ToLower(chi.URLParam(r, "type"))
	if metricType != gauge && metricType != counter {
		http.Error(w, "Only metricType  gauge and counter in request are allowed!", http.StatusNotImplemented)
		return
	}
	metricName := strings.ToLower(chi.URLParam(r, "name"))
	metricValue := chi.URLParam(r, "value")

	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()
	logrus.Info("SaveMetric Value: ", metricType, metricName, metricValue)
	err := Storage.UpdateMetricByParameters(requestContext, metricName, metricType, metricValue)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func getValueMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, getJSONError("Only application/json  can be Content-Type"), http.StatusUnsupportedMediaType)
	}
	defer r.Body.Close()

	var m storage.Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		logrus.Error(err)
		http.Error(w, getJSONError("Fail on parse request"), http.StatusBadRequest)
		return
	}

	if m.ID == "" || m.MType == "" {
		http.Error(w, getJSONError("ID or MType is empty"), http.StatusBadRequest)
		return
	}

	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()
	valueMetric, err := Storage.GetMetric(requestContext, m.ID, m.MType)
	mar, _ := json.Marshal(valueMetric)
	logrus.Info("GetValueMetric Value: ", mar)
	if err == nil && m.Delta == nil && m.Value == nil {
		valueMetric.SetHash(SignKey)
		render.JSON(w, r, &valueMetric)
	} else {
		http.Error(w, getJSONError("Data Not Found"), http.StatusNotFound)
	}

}

func getMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
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

	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()
	valueMetric, err := Storage.GetMetric(requestContext, metricName, metricType)
	mar, _ := json.Marshal(valueMetric)
	logrus.Info("GetMetric Value: ", mar)
	if err == nil && valueMetric.GetMetricType() == metricType {
		w.Write([]byte(valueMetric.GetValue()))
	} else {
		logrus.Error(err)
		http.Error(w, "Value NotFound!", http.StatusNotFound)
		return
	}
}

func getAllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()
	metrics, err := Storage.GetMetrics(requestContext)
	if err != nil {
		logrus.Error(err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{ \"error\" : \"error get values\"}"))
		return
	}
	metricsString := concatenationMetrics(metrics)
	logrus.Info("Data received: " + metricsString)
	w.Write([]byte(metricsString))

	logrus.Info(r.RequestURI)
}
func pingDataBase(w http.ResponseWriter, r *http.Request) {
	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()
	err := Storage.Ping(requestContext)
	if err != nil {
		logrus.Error(err)
		http.Error(w, getJSONError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func concatenationMetrics(metrics map[string]*storage.Metric) string {
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
	logrus.Error(errorText)
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
