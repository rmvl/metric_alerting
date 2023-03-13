package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	storageRepository "yalerting/cmd/storage"
)

func FlushMetrics(storage storageRepository.StorageRepository, cfg ServerConfig) {
	flusherIntervalTicker := time.NewTicker(time.Duration(cfg.GetStoreInterval()) * time.Second)

	fileName := cfg.StoreFile

	for {
		for range flusherIntervalTicker.C {
			producer, err := NewProducer(fileName)
			if err != nil {
				log.Fatal(err)
			}

			for k, v := range storage.GetCounters() {
				event := Metrics{ID: k, Delta: &v, MType: "counter"}
				if err := producer.WriteEvent(event); err != nil {
					log.Fatal(err)
				}
			}

			for k, v := range storage.GetGaugeMetrics() {
				floatVal, _ := strconv.ParseFloat(v, 64)
				if err != nil {
					fmt.Println(err)
					return
				}

				event := Metrics{ID: k, Value: &floatVal, MType: "gauge"}
				if err := producer.WriteEvent(event); err != nil {
					log.Fatal(err)
				}
			}

			producer.Close()
		}
	}
}

func RestoreMetrics(storage storageRepository.StorageRepository, cfg ServerConfig) {
	fileName := cfg.StoreFile
	consumer, err := NewConsumer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	for {
		metric, err := consumer.ReadEvent()
		if err != nil {
			fmt.Println("metric read error", err)
			break
		}

		switch metric.MType {
		case "counter":
			storage.IncrementCounter(metric.ID, *metric.Delta)
		case "gauge":
			storage.SetGaugeMetric(metric.ID, strconv.FormatFloat(*metric.Value, 'g', -1, 64))
		default:
			fmt.Println("not supported metric type")
			return
		}
	}
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	err = file.Truncate(0)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteEvent(event Metrics) error {
	return p.encoder.Encode(&event)
}

func (p *producer) Close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) ReadEvent() (*Metrics, error) {
	event := &Metrics{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
