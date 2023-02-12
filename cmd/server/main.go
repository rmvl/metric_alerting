package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"yalerting/cmd/handlers"
	storageClient "yalerting/cmd/storage"
)

func main() {
	storage := storageClient.NewMemStorage()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricName}/{metricType}/{metricValue}", handlers.UpdateMetric(storage))
	})

	r.Get("/", handlers.MetricList(storage))

	r.Get("/value/{metricType}/{metricName}", handlers.Get(storage))

	// запуск сервера с адресом localhost, порт 8080
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
