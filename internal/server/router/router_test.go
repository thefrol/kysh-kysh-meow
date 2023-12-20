package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/metricas"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/scan"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router/httpio"
	"github.com/thefrol/kysh-kysh-meow/internal/server/storage"
)

func Test_MeowRouter(t *testing.T) {

	type testResponse struct {
		code        int
		ContentType string
		body        string // сейчас не используется
	}

	tests := []struct {
		name     string
		prePosts []string // urls to call before the main test
		method   string
		body     string
		response testResponse
		route    string
	}{
		{
			name:     "positive counter #1",
			prePosts: []string{"/update/counter/test1/100"},
			method:   http.MethodGet,
			route:    "/value/counter/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100",
			},
		},
		{
			name: "positive counter #2",
			prePosts: []string{"/update/counter/test1/100",
				"/update/counter/test1/100"},
			method: http.MethodGet,
			route:  "/value/counter/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "200",
			},
		},
		{
			name: "positive counter #3",
			prePosts: []string{"/update/counter/test1/100",
				"/update/counter/test1/100.1"}, // float значение сервер не пример
			method: http.MethodGet,
			route:  "/value/counter/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100",
			},
		},
		{
			name:     "positive gauge #1",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodGet,
			route:    "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100",
			},
		},
		{
			name: "positive gauge #2",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/100.1"},
			method: http.MethodGet,
			route:  "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100.1",
			},
		},
		{
			name: "positive gauge #3",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/100.1",
				"/update/gauge/test1/none"}, //это значение не примет
			method: http.MethodGet,
			route:  "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100.1",
			},
		},
		{
			name: "positive gauge #4",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/100.1",
				"/update/gauge/test1/2e0"},
			method: http.MethodGet,
			route:  "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "2",
			},
		},
		{
			name:     "positive gauge #5",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodGet,
			route:    "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100",
			},
		},
		{
			name:     "json value positive #1",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodPost,
			route:    "/value",
			body:     `{"type":"gauge","id":"test1"}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"gauge","id":"test1","value":100}`,
			},
		},
		{
			name: "json value positive #2",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/200"},
			method: http.MethodPost,
			route:  "/value",
			body:   `{"type":"gauge","id":"test1"}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"gauge","id":"test1","value":200}`,
			},
		},
		{
			name: "json value positive #3",
			prePosts: []string{"/update/counter/test1/100",
				"/update/counter/test1/100"},
			method: http.MethodPost,
			route:  "/value",
			body:   `{"type":"counter","id":"test1"}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"counter","id":"test1","delta":200}`,
			},
		},

		{
			name:     "json value negative #1",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodGet,
			route:    "/value",
			body:     `{"type":"gauge","id":"test1"}`,
			response: testResponse{
				code:        http.StatusNotFound,
				ContentType: "text/plain",
				body:        "", // if == "" dont check
			},
		},
		{
			name:     "json value negative #2",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodPost,
			route:    "/value",
			body:     `{"type":"gauge","id":"test2"}`,
			response: testResponse{
				code:        http.StatusNotFound,
				ContentType: "",
				body:        "",
			},
		},
		{
			name:     "json value negative #3 wrong type",
			prePosts: []string{},
			method:   http.MethodPost,
			route:    "/value",
			body:     `{"type":"unknown","id":"test2"}`,
			response: testResponse{
				code:        http.StatusBadRequest,
				ContentType: "",
				body:        "",
			},
		},

		{
			name: "json update positive #1",
			prePosts: []string{"/update/counter/test1/100",
				"/update/counter/test1/100"},
			method: http.MethodPost,
			route:  "/update",
			body:   `{"type":"counter","id":"test1","delta":100}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"counter","id":"test1","delta":300}`,
			},
		},
		{
			name: "json update positive #2",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/101.2"},
			method: http.MethodPost,
			route:  "/update",
			body:   `{"type":"gauge","id":"test1","value":100.1}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"gauge","id":"test1","value":100.1}`,
			},
		},
		{
			name: "json update negative #1 unknown metric",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/101.2"},
			method: http.MethodPost,
			route:  "/update",
			body:   `{"type":"unknown","id":"test1","gauge":100.1}`,
			response: testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "json update negative #2 crazy json",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/101.2"},
			method: http.MethodPost,
			route:  "/update",
			body:   `{"type":"gauge","id":"test1","gauge":100.1}`,
			response: testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "json update negative #3 not a json",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/101.2"},
			method: http.MethodPost,
			route:  "/update",
			body:   `lalalal`,
			response: testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:     "list all metrics",
			prePosts: []string{},
			method:   http.MethodGet,
			route:    "/",
			body:     "",
			response: testResponse{
				code: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// создадим сервер
			router := NewMemRouter()
			server := httptest.NewServer(router)
			defer server.Close()

			// создадим клиента
			client := resty.New()

			// сделаем подготовительные запросы, которые
			// приведут сервер в правильное состояние
			for _, u := range tt.prePosts {
				_, err := client.R().SetHeader(httpio.HeaderContentType, "text/plain").Post(server.URL + u)
				assert.NoError(t, err, "error on preparing data for final request")
			}

			// выполним тестирующий запрос
			resp, err := client.R().SetBody(tt.body).Execute(tt.method, server.URL+tt.route)
			require.NoErrorf(t, err, "error on final request %+v", err)

			// проверим вывод
			assert.Equal(t, tt.response.code, resp.StatusCode())
			if tt.response.ContentType != "" {
				assert.Contains(t, resp.Header().Get(httpio.HeaderContentType), tt.response.ContentType)
			}

			if tt.response.body != "" {
				assert.JSONEq(t, tt.response.body, string(resp.Body()))
			}

		})
	}

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
			// подготовим сервер и запрос
			r := httptest.NewRequest(tt.method, tt.route, nil)
			w := httptest.NewRecorder()

			// выполним запрос
			NewMemRouter().ServeHTTP(w, r)

			// получим ответ
			result := w.Result()
			defer result.Body.Close()

			// проверим ответ
			assert.Equal(t, tt.response.code, result.StatusCode)
			assert.Contains(t, result.Header.Get(httpio.HeaderContentType), tt.response.ContentType)
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
			r := httptest.NewRequest(tt.method, tt.route, nil)
			w := httptest.NewRecorder()

			// запустим запрос
			NewMemRouter().ServeHTTP(w, r)

			// обработаем результат
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.response.code, result.StatusCode)
			assert.Contains(t, result.Header.Get(httpio.HeaderContentType), tt.response.ContentType)
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
			// подготовим тест-окружение
			r := httptest.NewRequest(tt.method, tt.route, nil)
			w := httptest.NewRecorder()

			// запустим запрос
			NewMemRouter().ServeHTTP(w, r)

			// получим ответ
			result := w.Result()
			defer result.Body.Close()

			// проверим результат
			assert.Equal(t, tt.response.code, result.StatusCode)
			assert.Contains(t, result.Header.Get(httpio.HeaderContentType), tt.response.ContentType)
			require.Equal(t, "", tt.response.body, "Тесты содержащие Body сейчас не поддерживаеются")
		})
	}

}

func NewMemRouter() http.Handler {
	s := storage.AsOperator(storage.New())

	// приготовим приложение
	// готовим репозитории
	counters := storage.CounterAdapter{
		Op: s,
	}

	gauges := storage.GaugeAdapter{
		Op: s,
	}

	// готовим прикладной уровень
	labels := scan.Labels{
		Labels: &storage.LabelsAdapter{Op: s},
	}

	reg := manager.Registry{
		Counters: &counters,
		Gauges:   &gauges,
	}

	man := metricas.Manager{
		Registry: reg,
	}

	// создаем роутер
	r := API{
		Manager:   man,
		Registry:  reg,
		Dashboard: labels,
	}

	return r.MeowRouter()
}
