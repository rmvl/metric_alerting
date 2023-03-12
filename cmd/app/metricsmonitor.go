package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const reportInterval = 10
const pollInterval = 2

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
}

func sendMetric(client http.Client, metric Metrics) error {
	body, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/update/", bytes.NewBuffer(body))
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
	fmt.Println(*metricResp.Value)
	resp.Body.Close()
	return nil
}

func MonitorMetrics() {
	//var trackedMetrics []string
	//var metricVal float64
	//trackedMetrics = []string{
	//	"Alloc",
	//	"BuckHashSys",
	//	"Frees",
	//	"GCCPUFraction",
	//	"GCSys",
	//	"HeapAlloc",
	//	"HeapIdle",
	//	"HeapInuse",
	//	"HeapObjects",
	//	"HeapReleased",
	//	"HeapSys",
	//	"LastGC",
	//	"Lookups",
	//	"MCacheInuse",
	//	"MCacheSys",
	//	"MSpanInuse",
	//	"MSpanSys",
	//	"Mallocs",
	//	"NextGC",
	//	"NumForcedGC",
	//	"NumGC",
	//	"OtherSys",
	//	"PauseTotalNs",
	//	"StackInuse",
	//	"StackSys",
	//	"Sys",
	//	"TotalAlloc",
	//	"PollCount",
	//	"RandomValue",
	//}
	client := http.Client{}

	metrics := MetricsToMonitor{}

	start := time.Now()
	pollTicker := time.NewTicker(pollInterval * time.Second)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval * time.Second)
	defer reportTicker.Stop()

	for {
		select {
		case x := <-reportTicker.C:
			fmt.Println(int(x.Sub(start).Seconds()))

			m := 494691.516772
			metricToSend := Metrics{
				ID:    "NextGC",
				MType: typeGauge,
				Value: &m,
			}
			if err := sendMetric(client, metricToSend); err != nil {
				fmt.Println(err)
			}

			//
			//fmt.Println(int(x.Sub(start).Seconds()))
			//for _, metric := range trackedMetrics {
			//	metricValue := reflect.Indirect(reflect.ValueOf(metrics)).FieldByName(metric)
			//
			//	if metricValue.CanUint() {
			//		metricVal = float64(metricValue.Uint())
			//	}
			//	if metricValue.CanInt() {
			//		metricVal = float64(metricValue.Int())
			//	}
			//	if metricValue.CanFloat() {
			//		metricVal = metricValue.Float()
			//	}
			//
			//	metricToSend := Metrics{
			//		ID:    metric,
			//		MType: typeGauge,
			//		Value: &metricVal,
			//	}
			//
			//	if err := sendMetric(client, metricToSend); err != nil {
			//		fmt.Println(err)
			//	}
			//}
			//
			//pollCountMetric := Metrics{
			//	ID:    "PollCount",
			//	MType: typeCounter,
			//	Delta: &metrics.PollCount,
			//}
			//if err := sendMetric(client, pollCountMetric); err != nil {
			//	fmt.Println(err)
			//}
			//
			//randomValueMetric := Metrics{
			//	ID:    "RandomValue",
			//	MType: typeCounter,
			//	Delta: &metrics.RandomValue,
			//}
			//if err := sendMetric(client, randomValueMetric); err != nil {
			//	fmt.Println(err)
			//}

			metrics.PollCount = 0
		case y := <-pollTicker.C:
			fmt.Println(int(y.Sub(start).Seconds()))
			runtime.ReadMemStats(&metrics.MemStats)
			metrics.PollCount += 1
			metrics.RandomValue = rand.Int63()
		}
	}

}
