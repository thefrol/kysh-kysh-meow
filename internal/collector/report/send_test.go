package report_test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/report"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

func TestSend(t *testing.T) {
	tests := []struct {
		name     string
		metricas []metrica.Metrica
		wantErr  bool
	}{
		{
			name: "counters",
			metricas: []metrica.Metrica{
				{ID: "test1", Delta: wrapInt64(10), MType: "counter"},
				{ID: "Test1", Delta: wrapInt64(9), MType: "counter"}},
			wantErr: false,
		},
		{
			name: "gauges",
			metricas: []metrica.Metrica{
				{ID: "test1", Value: wrapFloat64(10), MType: "gauge"},
				{ID: "Test1", Value: wrapFloat64(9), MType: "gauge"}},
			wantErr: false,
		},
		{
			name: "unknown metricas",
			metricas: []metrica.Metrica{
				{ID: "test1", Value: wrapFloat64(10), MType: "unknown"},
				{ID: "Test22", Value: wrapFloat64(9), MType: "unknown"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testHandler{}
			s := httptest.NewServer(&h)

			err := report.Send(tt.metricas, s.URL)
			if tt.wantErr {
				assert.Error(t, err, "Должна быть ошибка")
			}

			require.Equal(t, len(h.requests), 1, "Должен быть только один запрос на сервер")
			req := h.requests[0]
			defer req.Body.Close()

			sended := make([]metrica.Metrica, 0, len(tt.metricas))
			body, err := req.GetBody()
			require.NoError(t, err, "Ошибка получения тела")
			err = json.NewDecoder(body).Decode(&sended)
			require.NoError(t, err, "Не возможно размаршалить джейсон из отправленного джейсона: %v", err)

			eq := reflect.DeepEqual(tt.metricas, sended)
			assert.True(t, eq, "Ожидаемые к отправке данные не совпадают с полученными сервером")
		})
	}
}

func wrapInt64(v int) (ref *int64) {
	ref = new(int64)
	*ref = int64(v)
	return
}

func wrapFloat64(v int) (ref *float64) {
	ref = new(float64)
	*ref = float64(v)
	return
}

// testHandler простая структура, которая запоминает по каким адресам к ней делали запросы
// используется в тестировании.
type testHandler struct {
	requests []*http.Request
}

// ServerHTTP отвечает на запросы к серверы, исполяет интерфейс http.Handle
func (server *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var bb []byte
	defer r.Body.Close()

	// разархивируем если надо
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		unzipped, _ := gzip.NewReader(r.Body)
		defer unzipped.Close()
		bb, _ = io.ReadAll(unzipped)

	} else {
		bb, _ = io.ReadAll(r.Body)

	}

	// чтобы потом можно было прочитать тело запроса,
	// мы его сохраняем в сцециализированную переменную внутри запроса,
	// которая работает как замыкаение
	r.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(bb)), nil
	}

	// добавляем запросы в массив запросов. Теперь каждый такой запрос помнит и тело своего запроса
	server.requests = append(server.requests, r)
}

// routesUsed возвращает марштуры по которым проходили запросы
func (server testHandler) routesUsed() (routes []string) {
	for _, r := range server.requests {
		routes = append(routes, r.URL.Path)
	}
	return
}

// containsRoute возвращает true, если в принятых сервером маршрутах находится такой.
// Задается регулярным вырежаением pattern
func (server testHandler) containsRoute(pattern string) (bool, error) {
	for _, r := range server.requests {
		found, err := regexp.MatchString(pattern, r.URL.Path)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

// возвращает количество полученных запросов
func (server testHandler) NumRequests() int {
	return len(server.requests)
}

// stringerWrap позволяет использовать интерфейс стрингер в полях структуры. Позволяя запихивать туда люую переменную, которую можно обратить в строку
type stringerWrap struct {
	text string
}

func (s stringerWrap) String() string {
	return s.text
}
