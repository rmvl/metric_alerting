package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

func HashMetric(metric Metrics, secretKey string) (string, error) {
	var dataToHash string
	switch metric.MType {
	case "counter":
		//if metric.Delta == nil {
		//
		//}
		dataToHash = fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta)
	case "gauge":
		dataToHash = fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)
	default:
		fmt.Println("error")
		return "", errors.New("metric type is not implemented")
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(dataToHash))
	dst := h.Sum(nil)

	r := base64.URLEncoding.EncodeToString(dst)

	return r, nil
}

func CheckHash(metric Metrics, secretKey string) error {
	fmt.Println("merickhas", metric.Hash, &metric.Hash)
	if len(metric.Hash) == 0 {
		return errors.New("hash is not valid")
	}

	hash, err := HashMetric(metric, secretKey)
	if err != nil || hash != metric.Hash {
		return errors.New("hash is not valid")
	}

	return nil
}
