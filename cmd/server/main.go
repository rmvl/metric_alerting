package main

import (
	"fmt"
	"net/http"
	"yalerting/cmd/handlers"
	storageClient "yalerting/cmd/storage"
)

func main() {
	storage := storageClient.NewMemStorage()

	// маршрутизация запросов обработчику
	http.HandleFunc("/", handlers.HandleMetric(storage))
	// запуск сервера с адресом localhost, порт 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
