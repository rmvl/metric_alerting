package storage

import (
	"strconv"
	"sync"
)

type StorageRepository interface {
	SetGaugeMetric(name string, value string)
	IncrementCounter(name string, value int64)
	GetList() map[string]string
	GetCounterMetric(metricName string) (int64, bool)
	GetGaugeMetric(metricName string) (string, bool)
	GetMetric(metricName string, metricType string) (string, bool)
}

type MemStorage struct {
	metrics  map[string]string
	counters map[string]int64
	mutex    sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		make(map[string]string, 100),
		make(map[string]int64, 100),
		sync.RWMutex{},
	}
}

func (storage *MemStorage) SetGaugeMetric(name string, value string) {
	storage.mutex.Lock()
	storage.metrics[name] = value
	storage.mutex.Unlock()
}

func (storage *MemStorage) IncrementCounter(name string, value int64) {
	storage.mutex.Lock()

	_, ok := storage.counters[name]
	if !ok {
		storage.counters[name] = value
	} else {
		storage.counters[name] += value
	}
	storage.mutex.Unlock()
}

func (storage *MemStorage) GetCounterMetric(name string) (int64, bool) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	val, ok := storage.counters[name]
	if !ok {
		return 0, false
	}

	return val, true
}

func (storage *MemStorage) GetGaugeMetric(name string) (string, bool) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	val, ok := storage.metrics[name]
	if !ok {
		return "", false
	}

	return val, true
}

func (storage *MemStorage) GetMetric(name string, metricType string) (string, bool) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if metricType == "counter" {
		val, ok := storage.counters[name]
		if !ok {
			return "", false
		}

		return strconv.FormatInt(val, 10), true
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

func (storage *MemStorage) GetList() map[string]string {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	allMetrics := make(map[string]string, 100)
	for k, v := range storage.metrics {
		allMetrics[k] = v
	}
	for k, v := range storage.counters {
		allMetrics[k] = strconv.FormatInt(v, 10)
	}

	return allMetrics
}
