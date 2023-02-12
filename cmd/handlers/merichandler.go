package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	storageRepository "yalerting/cmd/storage"
)

func UpdateMetric(storage storageRepository.StorageRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "metricName")
		metricType := chi.URLParam(r, "metricType")
		metricValue := chi.URLParam(r, "metricValue")
		fmt.Println(metricName, metricType, metricValue)

		rw.Header().Set("Content-Type", "application/json")

		if metricType == "counter" {
			if s, err := strconv.ParseUint(metricType, 10, 32); err == nil {
				storage.IncrementCounter(metricName, s)
			}
		}
		if metricType == "gauge" {
			storage.SetGaugeMetric(metricName, metricValue)
		}
		fmt.Println(storage)

		rw.WriteHeader(http.StatusOK)
	}
}

func MetricList(storage storageRepository.StorageRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		rw.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(storage.GetList())
		fmt.Println(resp)
		rw.Write(resp)
	}
}

func Get(storage storageRepository.StorageRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "metricName")
		metricType := chi.URLParam(r, "metricType")
		fmt.Println(metricName, metricType)

		rw.Header().Set("Content-Type", "application/json")

		metrivVal, ok := storage.GetMetric(metricName, metricType)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(""))
			return
		}

		resp, _ := json.Marshal(metrivVal)

		rw.WriteHeader(http.StatusOK)
		rw.Write(resp)
	}
}
