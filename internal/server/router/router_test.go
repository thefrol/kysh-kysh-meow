package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thefrol/kysh-kysh-meow/internal/server/handlers"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

func init() {
	handlers.SetStore(storage.New())

}

// что такое test main?
func Test_updateCounter(t *testing.T) {
	type testResponse struct {
		code        int
		ContentType string
		body        string // сейчас не используется
	}

	tests := []struct {
		name     string
		method   string
		response testResponse
		route    string
	}{
		{
			name:   "positive counter #1",
			method: http.MethodPost,
			route:  "/update/counter/main/20",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "positive counter #2",
			method: http.MethodPost,
			route:  "/update/counter/main/-20",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "negative counter #1 GET",
			method: http.MethodGet,
			route:  "/update/counter/main/20",
			response: testResponse{
				code:        http.StatusNotFound,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "negative counter #2",
			method: http.MethodPost,
			route:  "/update/counter/main/20.1",
			response: testResponse{
				code:        http.StatusBadRequest,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "negative counter #3",
			method: http.MethodPost,
			route:  "/update/counter/main/none",
			response: testResponse{
				code:        http.StatusBadRequest,
				ContentType: "text/plain",
				body:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest(tt.method, tt.route, nil)
			w := httptest.NewRecorder()
			MeowRouter().ServeHTTP(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.response.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.response.ContentType)
			require.Equal(t, "", tt.response.body, "Тесты содержащие Body сейчас не поддерживаеются")
		})
	}

}

func Test_updateGauge(t *testing.T) {
	type testResponse struct {
		code        int
		ContentType string
		body        string // сейчас не используется
	}

	tests := []struct {
		name     string
		method   string
		response testResponse
		route    string
	}{

		{
			name:   "positive gauge #0",
			method: http.MethodPost,
			route:  "/update/gauge/main/20.11111111",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "positive gauge #1",
			method: http.MethodPost,
			route:  "/update/gauge/main/-20.11111111",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "positive gauge #2",
			method: http.MethodPost,
			route:  "/update/gauge/main/-20.1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "positive gauge #3",
			method: http.MethodPost,
			route:  "/update/gauge/main/-20",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "positive gauge #4",
			method: http.MethodPost,
			route:  "/update/gauge/main/2.11e13",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "positive gauge #5",
			method: http.MethodPost,
			route:  "/update/gauge/main/-0.11e23",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "positive gauge #6",
			method: http.MethodPost,
			route:  "/update/gauge/main/-0.11e-23",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "",
			},
		},

		{
			name:   "negative gauge #1",
			method: http.MethodPost,
			route:  "/update/gauge/main/none",
			response: testResponse{
				code:        http.StatusBadRequest,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "negative gauge #2 GET",
			method: http.MethodGet,
			route:  "/update/gauge/main/20",
			response: testResponse{
				code:        http.StatusNotFound,
				ContentType: "text/plain",
				body:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.SetStore(storage.New())
			r := httptest.NewRequest(tt.method, tt.route, nil)
			w := httptest.NewRecorder()
			MeowRouter().ServeHTTP(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.response.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.response.ContentType)
			require.Equal(t, "", tt.response.body, "Тесты содержащие Body сейчас не поддерживаеются")
		})
	}

}

func Test_updateUnknownType(t *testing.T) {
	type testResponse struct {
		code        int
		ContentType string
		body        string // сейчас не используется
	}

	tests := []struct {
		name     string
		method   string
		response testResponse
		route    string
	}{
		{
			name:   "negative unknown type #1",
			method: http.MethodPost,
			route:  "/update/karamba/main/none",
			response: testResponse{
				code:        http.StatusBadRequest,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "negative unknown type #2",
			method: http.MethodPost,
			route:  "/update/karamba/main/20",
			response: testResponse{
				code:        http.StatusBadRequest,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "negative unknown type #3",
			method: http.MethodPost,
			route:  "/update/karamba/main/2e10",
			response: testResponse{
				code:        http.StatusBadRequest,
				ContentType: "text/plain",
				body:        "",
			},
		},
		{
			name:   "negative unknown type #2 GET",
			method: http.MethodGet,
			route:  "/update/unknown/main/20",
			response: testResponse{
				code:        http.StatusNotFound,
				ContentType: "text/plain",
				body:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.route, nil)
			w := httptest.NewRecorder()
			MeowRouter().ServeHTTP(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.response.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.response.ContentType)
			require.Equal(t, "", tt.response.body, "Тесты содержащие Body сейчас не поддерживаеются")
		})
	}

}
