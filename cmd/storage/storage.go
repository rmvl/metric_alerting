package storage

type StorageRepository interface {
	SetGaugeMetric(name string, value string)
	IncrementCounter(name string, value uint64)
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
