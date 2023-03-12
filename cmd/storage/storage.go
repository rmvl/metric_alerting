package storage

import (
	"strconv"
	"sync"
)

type StorageRepository interface {
	SetGaugeMetric(name string, value string, withMutex bool)
	IncrementCounter(name string, value int64, withMutex bool)
	GetList() map[string]string
	GetCounters() map[string]int64
	GetGaugeMetrics() map[string]string
	GetCounterMetric(metricName string) (int64, bool)
	GetGaugeMetric(metricName string) (string, bool)
	GetMetric(metricName string, metricType string) (string, bool)
	GetMutex() *sync.RWMutex
}

type MemStorage struct {
	metrics  map[string]string
	counters map[string]int64
	Mutex    sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		make(map[string]string, 100),
		make(map[string]int64, 100),
		sync.RWMutex{},
	}
}

func (storage *MemStorage) GetCounters() map[string]int64 {
	return storage.counters
}

func (storage *MemStorage) GetGaugeMetrics() map[string]string {
	return storage.metrics
}

func (storage *MemStorage) SetGaugeMetric(name string, value string, withMutex bool) {
	if withMutex {
		storage.Mutex.Lock()
		storage.Mutex.Unlock()
	}
	storage.metrics[name] = value
}

func (storage *MemStorage) IncrementCounter(name string, value int64, withMutex bool) {
	if withMutex {
		storage.Mutex.Lock()
		storage.Mutex.Unlock()
	}

	_, ok := storage.counters[name]
	if !ok {
		storage.counters[name] = value
	} else {
		storage.counters[name] += value
	}
}

func (storage *MemStorage) GetCounterMetric(name string) (int64, bool) {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	val, ok := storage.counters[name]
	if !ok {
		return 0, false
	}

	return val, true
}

func (storage *MemStorage) GetGaugeMetric(name string) (string, bool) {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	val, ok := storage.metrics[name]
	if !ok {
		return "", false
	}

	return val, true
}

func (storage *MemStorage) GetMutex() *sync.RWMutex {
	return &storage.Mutex
}

func (storage *MemStorage) GetMetric(name string, metricType string) (string, bool) {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

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
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	allMetrics := make(map[string]string, 100)
	for k, v := range storage.metrics {
		allMetrics[k] = v
	}
	for k, v := range storage.counters {
		allMetrics[k] = strconv.FormatInt(v, 10)
	}

	return allMetrics
}
