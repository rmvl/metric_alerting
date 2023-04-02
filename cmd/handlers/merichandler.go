package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"yalerting/cmd/app"
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
			http.Error(rw, "metricType param is invalid", http.StatusNotImplemented)
			return
		}

		switch metricType {
		case "counter":
			if s, err := strconv.ParseInt(metricValue, 10, 64); err == nil {
				storage.IncrementCounter(metricName, s)
			} else {
				http.Error(rw, "metricValue param is not int64", http.StatusBadRequest)
				return
			}
		case "gauge":
			if _, err := strconv.ParseFloat(metricValue, 64); err == nil {
				storage.SetGaugeMetric(metricName, metricValue)
			} else {
				http.Error(rw, "metricValue param is not float 64", http.StatusBadRequest)
				return
			}
		default:
			http.Error(rw, "Unsupported metricType"+metricType, http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}

func UpdateMetricByJSONData(storage storageRepository.StorageRepository, cfg app.ServerConfig) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var metric app.Metrics
		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			http.Error(rw, "Not valid json", http.StatusBadRequest)
			return
		}

		if len(cfg.Key) > 0 {
			err := app.CheckHash(metric, cfg.Key)
			if err != nil {
				http.Error(rw, "hash not valid", http.StatusBadRequest)
				return
			}
		}

		if metric.ID == "" {
			http.Error(rw, "metricId param is empty", http.StatusBadRequest)
			return
		}
		if metric.MType == "" {
			http.Error(rw, "metricType param is missed", http.StatusBadRequest)
			return
		}

		if metric.MType != "counter" && metric.MType != "gauge" {
			http.Error(rw, "metricType param is invalid", http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case "counter":
			storage.IncrementCounter(metric.ID, *metric.Delta)
		case "gauge":
			storage.SetGaugeMetric(metric.ID, strconv.FormatFloat(*metric.Value, 'g', -1, 64))
		default:
			http.Error(rw, "Unsupported metricType"+metric.MType, http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(metric)
		if err != nil {
			http.Error(rw, "Can not prepare answer", http.StatusBadRequest)
			return
		}
		rw.Write(body)
		rw.WriteHeader(http.StatusOK)
	}
}

func MetricList(storage storageRepository.StorageRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html; charset=UTF-8")

		rw.WriteHeader(http.StatusOK)
		var response string

		for k, v := range storage.GetList() {
			response += k + ": " + v + "<br/>"
		}

		rw.Write([]byte(response))
	}
}

func GetMetricInJSON(storage storageRepository.StorageRepository, cfg app.ServerConfig) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		var metric app.Metrics
		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			http.Error(rw, "Not valid json", http.StatusBadRequest)
			return
		}

		if len(cfg.Key) > 0 {
			err := app.CheckHash(metric, cfg.Key)
			if err != nil {
				http.Error(rw, "hash not valid", http.StatusBadRequest)
				return
			}
		}

		if metric.ID == "" {
			http.Error(rw, "metricId param is empty", http.StatusBadRequest)
			return
		}
		if metric.MType == "" {
			http.Error(rw, "metricType param is missed", http.StatusBadRequest)
			return
		}

		if metric.MType == "" || metric.ID == "" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(""))
			return
		}

		switch metric.MType {
		case "counter":
			val, ok := storage.GetCounterMetric(metric.ID)
			if !ok {
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte(""))
				return
			}
			metric.Delta = &val
		case "gauge":
			val, ok := storage.GetGaugeMetric(metric.ID)
			if !ok {
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte(""))
				return
			}
			gaugeMetric, err := strconv.ParseFloat(val, 64)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(""))
				return
			}
			metric.Value = &gaugeMetric
		default:
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("Not supported metric type" + metric.MType))
			return
		}

		if len(cfg.Key) > 0 {
			hash, err := app.HashMetric(metric, cfg.Key)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("failed to hash"))
			}
			metric.Hash = hash
		}

		rw.WriteHeader(http.StatusOK)
		body, err := json.Marshal(metric)
		if err != nil {
			http.Error(rw, "Can not prepare answer", http.StatusBadRequest)
			return
		}
		rw.Write(body)
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
		rw.Header().Set("Content-Type", "text/html")

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
