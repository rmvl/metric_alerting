package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

const reportInterval = 10
const pollInterval = 2

const typeGauge = "gauge"
const typeCounter = "counter"

type metrics struct {
	PollCount   int
	RandomValue int
	runtime.MemStats
}

func sendMetric(client http.Client, metricType string, metricName string, metricValue string) {
	var body = []byte(`{"message":"Hello"}`)
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/update/"+metricType+"/"+metricName+"/"+metricValue, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	_, errC := client.Do(request)
	if errC != nil {
		fmt.Println(errC)
	}
}

func main() {
	var trackedMetrics []string
	var metricVal string
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

	metrics := metrics{}

	start := time.Now()
	pollTicker := time.NewTicker(pollInterval * time.Second)
	reportTicker := time.NewTicker(reportInterval * time.Second)
	for {
		select {
		case x := <-reportTicker.C:
			fmt.Println(int(x.Sub(start).Seconds()))
		case y := <-pollTicker.C:
			fmt.Println(int(y.Sub(start).Seconds()))
			runtime.ReadMemStats(&metrics.MemStats)
			metrics.PollCount += 1
			metrics.RandomValue = rand.Int()
			for _, metric := range trackedMetrics {
				fmt.Println(metric)
				metricValue := reflect.Indirect(reflect.ValueOf(metrics)).FieldByName(metric)
				switch reflect.TypeOf(metricValue).String() {
				case "uint64":
					metricVal = strconv.FormatUint(metricValue.Uint(), 10)
				case "float64":
					metricVal = strconv.FormatFloat(metricValue.Float(), 'g', 5, 64)
				}

				sendMetric(client, metric, typeGauge, metricVal)
			}
			sendMetric(client, "PollCount", typeCounter, strconv.Itoa(metrics.PollCount))
			sendMetric(client, "RandomValue", typeGauge, strconv.Itoa(metrics.RandomValue))
		}
	}
}
