package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	storageClient "yalerting/cmd/storage"
)

func TestHandleMetric(t *testing.T) {
	storage := storageClient.NewMemStorage()

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "application/json",
				statusCode:  200,
			},
			request: "/update/PollCount/counter/1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(UpdateMetric(storage))
			h(w, request)
			result := w.Result()
			defer result.Body.Close()

			// проверяем код ответа
			fmt.Println(result)
			if result.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
			}

			fmt.Println(result)
		})
	}
}
