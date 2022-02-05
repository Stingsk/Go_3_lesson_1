package main

import (
	metrics "github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"log"
	"net/http"
)

func main() {
	go metrics.NewMonitor(1000 * 1000 * 10)
	http.HandleFunc("/GetMetrics/", metrics.GetMetrics)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
