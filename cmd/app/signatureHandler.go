package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
)

func HashMetric(metric *Metrics, secretKey string) (string, error) {
	var dataToHash string
	switch metric.MType {
	case "counter":
		dataToHash = fmt.Sprintf("%s:counter:%d", metric.ID, metric.Delta)
	case "gauge":
		dataToHash = fmt.Sprintf("%s:gauge:%d", metric.ID, metric.Value)
	default:
		fmt.Println("error")
		return "", errors.New("metric type is not implemented")
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(dataToHash))
	dst := h.Sum(nil)
	return string(dst), nil
}

func CheckHash(metric *Metrics, secretKey string) error {
	hash, err := HashMetric(metric, secretKey)
	if err != nil || hash != metric.Hash {
		return errors.New("hash is not valid")
	}

	return nil
}
