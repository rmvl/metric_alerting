package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yalerting/cmd/app"
	storageRepository "yalerting/cmd/storage"
)

func UpdateMetric(storage storageRepository.StorageRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var metric app.Metrics
		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			http.Error(rw, "Not valid json", http.StatusBadRequest)
			return
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
			http.Error(rw, "metricType param is invalid", http.StatusNotImplemented)
			return
		}

		switch metric.MType {
		case "counter":
			storage.IncrementCounter(metric.ID, *metric.Delta)
		case "gauge":
			storage.SetGaugeMetric(metric.ID, strconv.FormatFloat(*metric.Value, 'g', 5, 64))
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
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

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
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var metric app.Metrics
		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			http.Error(rw, "Not valid json", http.StatusBadRequest)
			return
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
		rw.Header().Set("Content-Type", "text/html")

		switch metric.MType {
		case "counter":
			val, ok := storage.GetCounterMetric(metric.ID)
			if ok != true {
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte(""))
				return
			}
			metric.Delta = &val
		case "gauge":
			val, ok := storage.GetGaugeMetric(metric.ID)
			if ok != true {
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

		rw.WriteHeader(http.StatusOK)
		body, err := json.Marshal(metric)
		if err != nil {
			http.Error(rw, "Can not prepare answer", http.StatusBadRequest)
			return
		}
		rw.Write(body)
	}
}
