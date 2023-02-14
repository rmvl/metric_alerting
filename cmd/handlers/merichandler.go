package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	storageRepository "yalerting/cmd/storage"
)

func HandleMetric(storage storageRepository.StorageRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		fmt.Println(r.URL.String())
		params := strings.Split(r.URL.String(), "/")
		if len(params) < 5 {
			http.Error(rw, "not enough params to save metric", http.StatusBadRequest)
			return
		}
		if len(params) > 5 {
			http.Error(rw, "too many params to save metric", http.StatusBadRequest)
			return
		}
		if params[3] == "counter" {
			if s, err := strconv.ParseUint(params[3], 10, 32); err == nil {
				storage.IncrementCounter(params[2], s)
			}
		}
		if params[3] == "gauge" {
			storage.SetGaugeMetric(params[2], params[4])
		}
		fmt.Println(storage)

		rw.WriteHeader(http.StatusOK)
	}
}
