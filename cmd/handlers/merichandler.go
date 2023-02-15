package handlers

import (
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

		if metricName == "" {
			http.Error(rw, "metricName param is missed", http.StatusBadRequest)
			return
		}
		if metricType == "" {
			http.Error(rw, "metricType param is missed", http.StatusBadRequest)
			return
		}
		if metricValue == "" {
			http.Error(rw, "metricValue param is missed", http.StatusBadRequest)
			return
		}

		if metricType != "counter" && metricType != "gauge" {
			http.Error(rw, "metricValue param is missed", http.StatusNotImplemented)
			return
		}

		if metricType == "counter" {
			if s, err := strconv.ParseUint(metricValue, 10, 64); err == nil {
				storage.IncrementCounter(metricName, s)
			} else {
				http.Error(rw, "metricValue param is not int64", http.StatusBadRequest)
				return
			}
		}
		if metricType == "gauge" {
			if _, err := strconv.ParseFloat(metricValue, 64); err == nil {
				storage.SetGaugeMetric(metricName, metricValue)
			} else {
				http.Error(rw, "metricValue param is not float 64", http.StatusBadRequest)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
	}
}

func MetricList(storage storageRepository.StorageRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")

		rw.WriteHeader(http.StatusOK)
		var response string

		for k, v := range storage.GetList() {
			response += k + ": " + v + "<br/>"
		}

		rw.Write([]byte(response))
	}
}

func GetMetric(storage storageRepository.StorageRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "metricName")
		metricType := chi.URLParam(r, "metricType")
		if metricType == "" {
			http.Error(rw, "metricType param is missed", http.StatusBadRequest)
			return
		}
		if metricName == "" {
			http.Error(rw, "metricName param is missed", http.StatusBadRequest)
			return
		}

		if metricType == "" || metricName == "" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(""))
			return
		}
		rw.Header().Set("Content-Type", "text/plain")

		metrivVal, ok := storage.GetMetric(metricName, metricType)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(""))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(metrivVal))
	}
}
