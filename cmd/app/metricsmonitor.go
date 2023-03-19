package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const typeGauge = "gauge"
const typeCounter = "counter"

type MetricsToMonitor struct {
	PollCount   int64
	RandomValue int64
	runtime.MemStats
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func sendMetric(client http.Client, metric Metrics, cfg AgentConfig) error {
	if len(cfg.Key) > 0 {
		hash, err := HashMetric(&metric, cfg.Key)
		if err != nil {
			fmt.Println(err)
			return err
		}
		metric.Hash = hash
	}

	body, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodPost, "http://"+cfg.Address+"/update/", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, errC := client.Do(request)
	if errC != nil {
		fmt.Println(errC)
		return errC
	}

	var metricResp Metrics
	errR := json.NewDecoder(resp.Body).Decode(&metricResp)
	if errR != nil {
		fmt.Println(errR)
	}
	resp.Body.Close()
	return nil
}

func MonitorMetrics(cfg AgentConfig) {
	var trackedMetrics []string
	var metricVal float64
	trackedMetrics = []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"PollCount",
		"RandomValue",
	}
	client := http.Client{}

	metrics := MetricsToMonitor{}

	start := time.Now()

	pollTicker := time.NewTicker(time.Duration(cfg.GetPollInterval()) * time.Second)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Duration(cfg.GetReportInterval()) * time.Second)
	defer reportTicker.Stop()

	for {
		select {
		case x := <-reportTicker.C:
			fmt.Println(int(x.Sub(start).Seconds()))
			for _, metric := range trackedMetrics {
				metricValue := reflect.Indirect(reflect.ValueOf(metrics)).FieldByName(metric)

				if metricValue.CanUint() {
					metricVal = float64(metricValue.Uint())
				}
				if metricValue.CanInt() {
					metricVal = float64(metricValue.Int())
				}
				if metricValue.CanFloat() {
					metricVal = metricValue.Float()
				}

				metricToSend := Metrics{
					ID:    metric,
					MType: typeGauge,
					Value: &metricVal,
				}

				if err := sendMetric(client, metricToSend, cfg); err != nil {
					fmt.Println(err)
				}
			}

			pollCountMetric := Metrics{
				ID:    "PollCount",
				MType: typeCounter,
				Delta: &metrics.PollCount,
			}
			if err := sendMetric(client, pollCountMetric, cfg); err != nil {
				fmt.Println(err)
			}

			randomValueMetric := Metrics{
				ID:    "RandomValue",
				MType: typeCounter,
				Delta: &metrics.RandomValue,
			}
			if err := sendMetric(client, randomValueMetric, cfg); err != nil {
				fmt.Println(err)
			}

			metrics.PollCount = 0
		case y := <-pollTicker.C:
			fmt.Println(int(y.Sub(start).Seconds()))
			runtime.ReadMemStats(&metrics.MemStats)
			metrics.PollCount += 1
			metrics.RandomValue = rand.Int63()
		}
	}

}
