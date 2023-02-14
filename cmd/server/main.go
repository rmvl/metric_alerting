package main

import (
	"fmt"
	"net/http"
)

func HandleMetric(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("store or smth"))
	//fmt.Println(r)
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", HandleMetric)
	// запуск сервера с адресом localhost, порт 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
