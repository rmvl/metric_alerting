package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"yalerting/cmd/app"
	storageClient "yalerting/cmd/storage"
)

func TestUpdateMetric(t *testing.T) {
	storage := storageClient.NewMemStorage()

	type want struct {
		contentType string
		statusCode  int
		metricResp  app.Metrics
	}
	var delta int64
	var gaugeValue float64
	delta = 200
	gaugeValue = 300.123
	tests := []struct {
		name      string
		metricReq app.Metrics
		want      want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "application/json",
				statusCode:  200,
				metricResp:  app.Metrics{ID: "PollCount", MType: "counter", Delta: &delta},
			},
			metricReq: app.Metrics{ID: "PollCount", MType: "counter", Delta: &delta},
		},
		{
			name: "simple test #1",
			want: want{
				contentType: "application/json",
				statusCode:  200,
				metricResp:  app.Metrics{ID: "Alloc", MType: "gauge", Value: &gaugeValue},
			},
			metricReq: app.Metrics{ID: "Alloc", MType: "gauge", Value: &gaugeValue},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := "/update"
			body, _ := json.Marshal(tt.metricReq)

			request := httptest.NewRequest(http.MethodPost, method, bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			h := http.HandlerFunc(UpdateMetricByJSONData(storage))
			h(w, request)
			result := w.Result()
			defer result.Body.Close()

			// проверяем код ответа
			if result.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
			}

			var metricResp app.Metrics
			errR := json.NewDecoder(result.Body).Decode(&metricResp)
			if errR != nil {
				fmt.Println(errR)
			}

			if reflect.DeepEqual(metricResp, tt.want.metricResp) != true {
				wantResp, _ := json.Marshal(tt.want.metricResp)
				gotResp, _ := json.Marshal(metricResp)
				t.Errorf("expected resp %s, got %s", string(wantResp), string(gotResp))
			}
		})
	}
}

func TestMetricList(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		content     string
	}
	tests := []struct {
		name    string
		request string
		storage storageClient.StorageRepository
		want    want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "application/json",
				statusCode:  200,
				content:     "",
			},
			request: "/",
			storage: func() storageClient.StorageRepository {
				storage := storageClient.NewMemStorage()
				return storage
			}(),
		},
		{
			name: "simple test #2",
			want: want{
				contentType: "application/json",
				statusCode:  200,
				content:     "PollCount: 300<br/>",
			},
			request: "/",
			storage: func() storageClient.StorageRepository {
				storage := storageClient.NewMemStorage()
				storage.IncrementCounter("PollCount", 300)
				return storage
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(MetricList(tt.storage))
			h(w, request)
			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			data, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			metricList := string(data)

			require.NoError(t, err)

			// проверяем код ответа
			if result.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
			}
			if metricList != tt.want.content {
				t.Errorf("Expected response body %s, got %s", tt.want.content, metricList)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		metricResp  app.Metrics
	}

	var delta int64
	var gaugeValue float64
	delta = 200
	gaugeValue = 300.123

	tests := []struct {
		name      string
		storage   storageClient.StorageRepository
		metricReq app.Metrics
		want      want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "application/json",
				statusCode:  200,
				metricResp:  app.Metrics{ID: "PollCount", MType: "counter", Delta: &delta},
			},
			metricReq: app.Metrics{ID: "PollCount", MType: "counter"},
			storage: func() storageClient.StorageRepository {
				storage := storageClient.NewMemStorage()
				storage.IncrementCounter("PollCount", delta)
				return storage
			}(),
		},
		{
			name:      "simple test #2",
			metricReq: app.Metrics{ID: "Alloc", MType: "gauge"},
			want: want{
				contentType: "application/json",
				statusCode:  200,
				metricResp:  app.Metrics{ID: "Alloc", MType: "gauge", Value: &gaugeValue},
			},
			storage: func() storageClient.StorageRepository {
				storage := storageClient.NewMemStorage()
				storage.SetGaugeMetric("Alloc", "300.123")
				return storage
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := "/value/"
			body, _ := json.Marshal(tt.metricReq)
			request := httptest.NewRequest(http.MethodPost, method, bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			h := http.HandlerFunc(GetMetricInJSON(tt.storage))
			h(w, request)
			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			// проверяем код ответа
			if result.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
			}

			var metricResp app.Metrics
			errR := json.NewDecoder(result.Body).Decode(&metricResp)
			if errR != nil {
				fmt.Println(errR)
			}

			if reflect.DeepEqual(metricResp, tt.want.metricResp) != true {
				wantResp, _ := json.Marshal(tt.want.metricResp)
				gotResp, _ := json.Marshal(metricResp)
				t.Errorf("expected resp %s, got %s", string(wantResp), string(gotResp))
			}
		})
	}
}

func TestGetFail(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		metricResp  app.Metrics
	}

	tests := []struct {
		name      string
		storage   storageClient.StorageRepository
		metricReq app.Metrics
		want      want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "application/json",
				statusCode:  404,
			},
			metricReq: app.Metrics{ID: "PollCount", MType: "counter"},
			storage: func() storageClient.StorageRepository {
				storage := storageClient.NewMemStorage()
				return storage
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := "/value/"
			body, _ := json.Marshal(tt.metricReq)
			request := httptest.NewRequest(http.MethodPost, method, bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			h := http.HandlerFunc(GetMetricInJSON(tt.storage))
			h(w, request)
			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			// проверяем код ответа
			if result.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
			}
		})
	}
}
