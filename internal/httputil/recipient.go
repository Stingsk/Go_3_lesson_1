package httputil

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
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
			http.HandleFunc("/update/", recipient)
			log.Fatal(http.ListenAndServe(":8080", nil))
		}
	}
}

func recipient(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println(r.RequestURI)
}
