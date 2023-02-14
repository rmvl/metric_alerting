package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	storageClient "yalerting/cmd/storage"
)

func TestUpdateMetric(t *testing.T) {
	storage := storageClient.NewMemStorage()

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		request     string
		want        want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "application/json",
				statusCode:  200,
			},
			metricType:  "counter",
			metricName:  "PollCount",
			metricValue: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := "/update/" + tt.metricName + "/" + tt.metricType + "/" + tt.metricValue
			request := httptest.NewRequest(http.MethodPost, method, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", tt.metricType)
			rctx.URLParams.Add("metricName", tt.metricName)
			rctx.URLParams.Add("metricValue", tt.metricValue)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			h := http.HandlerFunc(UpdateMetric(storage))
			h(w, request)
			result := w.Result()
			defer result.Body.Close()

			// проверяем код ответа
			if result.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
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
		{
			name: "simple test #3",
			want: want{
				contentType: "application/json",
				statusCode:  200,
				content:     "PollCount: 300<br/>StackSys: 300.1<br/>",
			},
			request: "/",
			storage: func() storageClient.StorageRepository {
				storage := storageClient.NewMemStorage()
				storage.IncrementCounter("PollCount", 300)
				storage.SetGaugeMetric("StackSys", "300.1")
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
		content     string
	}
	tests := []struct {
		name        string
		metricType  string
		metricValue string
		storage     storageClient.StorageRepository
		want        want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "application/json",
				statusCode:  200,
				content:     "300",
			},
			metricType:  "counter",
			metricValue: "PollCount",
			storage: func() storageClient.StorageRepository {
				storage := storageClient.NewMemStorage()
				storage.IncrementCounter("PollCount", 300)
				return storage
			}(),
		},
		{
			name: "simple test #2",
			want: want{
				contentType: "application/json",
				statusCode:  404,
				content:     "",
			},
			metricType:  "counter",
			metricValue: "PollCount",
			storage: func() storageClient.StorageRepository {
				storage := storageClient.NewMemStorage()
				return storage
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := "/value/" + tt.metricType + "/" + tt.metricValue
			request := httptest.NewRequest(http.MethodGet, method, nil)

			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", "counter")
			rctx.URLParams.Add("metricName", "PollCount")

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			h := http.HandlerFunc(GetMetric(tt.storage))
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
