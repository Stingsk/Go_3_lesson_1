package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func setMetrics(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	// читаем Body
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println(b)
}

func main() {
	http.HandleFunc("/setMetrics/", setMetrics)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
