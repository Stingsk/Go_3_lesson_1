package main

import (
	"fmt"
	"log"
	"net/http"
)

func setMetrics(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println(r.RequestURI)
}

func main() {
	http.HandleFunc("/", setMetrics)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
