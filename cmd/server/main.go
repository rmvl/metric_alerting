package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"yalerting/cmd/handlers"
	storageClient "yalerting/cmd/storage"
)

type config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	storage := storageClient.NewMemStorage()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.MetricList(storage))
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", handlers.GetMetricInJSON(storage))
		r.Get("/{metricType}/{metricName}", handlers.GetMetric(storage))
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.UpdateMetricByJSONData(storage))
		r.Post("/{metricType}/{metricName}/{metricValue}", handlers.UpdateMetric(storage))
	})

	// запуск сервера с адресом localhost, порт 8080
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		fmt.Println(err)
	}
}
