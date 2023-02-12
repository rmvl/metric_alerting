package storage

import (
	"strconv"
)

type StorageRepository interface {
	SetGaugeMetric(name string, value string)
	IncrementCounter(name string, value uint64)
	GetList() map[string]string
	GetMetric(metricName string, metricType string) (string, bool)
}

type MemStorage struct {
	metrics  map[string]string
	counters map[string]uint64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		make(map[string]string, 100),
		make(map[string]uint64, 100),
	}
}

func (storage MemStorage) SetGaugeMetric(name string, value string) {
	storage.metrics[name] = value
}

func (storage MemStorage) IncrementCounter(name string, value uint64) {
	_, ok := storage.counters[name]
	if !ok {
		storage.counters[name] = value
	}
	storage.counters[name]++
}

func (storage MemStorage) GetMetric(name string, metricType string) (string, bool) {
	if metricType == "counter" {
		val, ok := storage.counters[name]
		if !ok {
			return "", false
		}

		return strconv.FormatUint(val, 10), true
	}

	if metricType == "gauge" {
		val, ok := storage.metrics[name]
		if !ok {
			return "", false
		}

		return val, true
	}

	return "", false
}

func (storage MemStorage) GetList() map[string]string {
	allMetrics := make(map[string]string, 100)
	for k, v := range storage.metrics {
		allMetrics[k] = v
	}
	for k, v := range storage.counters {
		allMetrics[k] = strconv.FormatUint(v, 10)
	}

	return allMetrics
}
