package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"yalerting/cmd/app"
	"yalerting/cmd/handlers"
	storageClient "yalerting/cmd/storage"
)

func main() {
	var cfg app.ServerConfig
	loadConfiguration(&cfg)

	storage := storageClient.NewMemStorage()

	if cfg.Restore {
		app.RestoreMetrics(storage, cfg)
	}
	//flush metrics to file
	go app.FlushMetrics(storage, cfg)

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
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		fmt.Println(err)
	}
}

func loadConfiguration(cfg *app.ServerConfig) {
	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}

	var address, storeInterval, storeFile string
	restore := new(app.Restore)
	flag.Var(restore, "r", "need to restore from file")
	flag.StringVar(&address, "a", cfg.Address, "server address")
	flag.StringVar(&storeInterval, "i", cfg.StoreInterval, "store interval")
	flag.StringVar(&storeFile, "f", cfg.StoreFile, "store file")
	flag.Parse()

	_, present := os.LookupEnv("ADDRESS")
	if !present && len(address) > 0 {
		cfg.Address = address
	}
	_, present = os.LookupEnv("RESTORE")
	if !present && !restore.IsSet {
		cfg.Restore = restore.Value
	}
	_, present = os.LookupEnv("STORE_INTERVAL")
	if !present && len(storeInterval) > 0 {
		cfg.StoreInterval = storeInterval
	}
	_, present = os.LookupEnv("STORE_FILE")
	if !present && len(storeFile) > 0 {
		cfg.StoreFile = storeFile
	}
}
