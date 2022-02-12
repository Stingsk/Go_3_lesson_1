package httputil

import (
	"context"
	"errors"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sync"
	"time"
)

func RunRecipient(ctx context.Context, wg *sync.WaitGroup, sigChan chan os.Signal) error {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sigChan:
			return errors.New("аварийное завершение")
		default:
			apiRouter := chi.NewRouter()
			setMiddlewares(apiRouter)
			apiRouter.Post("/update/*", recipient)

			logrus.Info("Starting HTTP server")

			err := http.ListenAndServe("localhost:8080", apiRouter)
			if err != nil {
				return err
			}
		}
	}
}

func recipient(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	logrus.Info(r.RequestURI)
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
