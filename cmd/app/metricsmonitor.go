package app

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

func sendMetric(client http.Client, metricType string, metricName string, metricValue string) error {
	var body = []byte("")
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/update/"+metricType+"/"+metricName+"/"+metricValue, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, errC := client.Do(request)
	if errC != nil {
		fmt.Println(errC)
		return errC
	}
	resp.Body.Close()
	return nil
}

func MonitorMetrics() {
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
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval * time.Second)
	defer reportTicker.Stop()

	for {
		select {
		case x := <-reportTicker.C:
			fmt.Println(int(x.Sub(start).Seconds()))
			for _, metric := range trackedMetrics {
				metricValue := reflect.Indirect(reflect.ValueOf(metrics)).FieldByName(metric)

				if metricValue.CanUint() {
					metricVal = strconv.FormatUint(metricValue.Uint(), 10)
				}
				if metricValue.CanInt() {
					metricVal = strconv.FormatInt(metricValue.Int(), 10)
				}
				if metricValue.CanFloat() {
					metricVal = strconv.FormatFloat(metricValue.Float(), 'g', 5, 64)
				}

				err := sendMetric(client, typeGauge, metric, metricVal)
				fmt.Println(err)
			}
			sendMetric(client, typeCounter, "PollCount", strconv.Itoa(metrics.PollCount))
			sendMetric(client, typeGauge, "RandomValue", strconv.Itoa(metrics.RandomValue))

			metrics.PollCount = 0
		case y := <-pollTicker.C:
			fmt.Println(int(y.Sub(start).Seconds()))
			runtime.ReadMemStats(&metrics.MemStats)
			metrics.PollCount += 1
			metrics.RandomValue = rand.Int()
		}
	}

}
